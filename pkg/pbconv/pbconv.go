package pbconv

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"time"

	"slices"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// defaultFields defines the default struct field names to be converted
// if processFields is not specified or is nil.
var defaultFields = []string{"CreatedAt", "UpdatedAt", "DeletedAt"}

// SliceStructTimeToPbTimestamp converts all time.Time fields specified in processFields
// from a slice of structs (fromObjSlice) to *timestamppb.Timestamp fields in the corresponding
// protobuf object slice (pbObjSlice).
//
// pbObjSlice and fromObjSlice must be pointers to slices of the same length.
// processFields is a pointer to a slice of field names to process; if nil or empty, defaultFields is used.
// diveFields specifies the field names to recursively process (for nested structs or slices).
//
// Example usage:
//
//	err := SliceStructTimeToPbTimestamp(&pbParkSlice, &parkSlice, nil, "Landmarks", "Event")
//
// This will convert pbParkSlice[i].Landmarks[j].CreatedAt, pbParkSlice[i].Landmarks[j].UpdatedAt,
// pbParkSlice[i].Event.DeletedAt, etc.
func SliceStructTimeToPbTimestamp(pbObjSlice any, fromObjSlice any, processFields *[]string, diveFields ...string) error {
	// Validate input types: both must be pointers to slices
	pbObjType := reflect.TypeOf(pbObjSlice)
	fromObjType := reflect.TypeOf(fromObjSlice)
	if pbObjType.Kind() != reflect.Ptr || pbObjType.Elem().Kind() != reflect.Slice ||
		fromObjType.Kind() != reflect.Ptr || fromObjType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("pbObjSlice and fromObjSlice must be pointers to slices")
	}

	// Extract slices
	pbObjSliceValue := reflect.ValueOf(pbObjSlice).Elem()
	fromObjSliceValue := reflect.ValueOf(fromObjSlice).Elem()

	// Ensure slices are the same length
	if pbObjSliceValue.Len() != fromObjSliceValue.Len() {
		return fmt.Errorf("pbObjSlice and fromObjSlice must be the same length")
	}

	var wg sync.WaitGroup
	// Iterate and process each element concurrently
	for i := range pbObjSliceValue.Len() {
		pbObj := pbObjSliceValue.Index(i).Addr().Interface()
		fromObj := fromObjSliceValue.Index(i).Addr().Interface()

		wg.Add(1)
		go func(pbObj, fromObj any) {
			defer wg.Done()
			StructTimeToPbTimestamp(pbObj, fromObj, processFields, diveFields...)
		}(pbObj, fromObj)
	}

	wg.Wait()
	return nil
}

// StructTimeToPbTimestamp converts all time.Time fields specified in processFields
// from a struct (fromObj) to *timestamppb.Timestamp fields in the corresponding
// protobuf object (pbObj).
//
// pbObj and fromObj must be pointers to structs with matching fields.
// processFields is a pointer to a slice of field names to process; if nil or empty, defaultFields is used.
// diveFields specifies the field names to recursively process (for nested structs or slices).
//
// Example usage:
//
//	err := StructTimeToPbTimestamp(&pbPark, &park, nil, "Landmarks", "Event")
//
// This will convert pbPark.Landmarks[i].CreatedAt, pbPark.Landmarks[i].UpdatedAt,
// pbPark.Event.DeletedAt, etc.
func StructTimeToPbTimestamp(pbObj any, fromObj any, processFields *[]string, diveFields ...string) error {
	fromValue := getStructValue(fromObj)
	pbValue := getStructValue(pbObj)

	for i := 0; i < pbValue.NumField(); i++ {
		pbFieldName := pbValue.Type().Field(i).Name

		// Recursively process nested fields if specified in diveFields
		if slices.Contains(diveFields, pbFieldName) {
			pbf := pbValue.Field(i)
			ff := fromValue.FieldByName(pbFieldName)

			if pbf.Kind() == reflect.Slice {
				for j := 0; j < pbf.Len(); j++ {
					err := StructTimeToPbTimestamp(pbf.Index(j).Addr().Interface(), ff.Index(j).Addr().Interface(), processFields, diveFields...)
					if err != nil {
						return err
					}
				}
			} else {
				if pbf.Kind() == reflect.Struct {
					err := StructTimeToPbTimestamp(pbf, ff, processFields, diveFields...)
					if err != nil {
						return err
					}
				}
				if pbf.Kind() == reflect.Ptr && !pbf.IsNil() {
					pbfp := pbf.Elem()
					if pbfp.Kind() != reflect.Struct {
						continue
					}
					err := StructTimeToPbTimestamp(pbfp.Addr().Interface(), ff.Addr().Interface(), processFields, diveFields...)
					if err != nil {
						return err
					}
				}
			}
			continue
		}

		// Determine which fields to process
		var fieldsToProcess []string
		if processFields != nil && len(*processFields) > 0 {
			fieldsToProcess = *processFields
		} else {
			fieldsToProcess = defaultFields
		}

		if !slices.Contains(fieldsToProcess, pbFieldName) {
			continue
		}

		fromField := fromValue.FieldByName(pbFieldName)
		pbField := pbValue.FieldByName(pbFieldName)

		if !fromField.IsValid() || !pbField.IsValid() {
			continue
		}

		// Process time.Time or *time.Time fields
		if fromField.Type() == reflect.TypeOf(time.Time{}) || fromField.Type() == reflect.TypeOf(&time.Time{}) {
			if fromField.Kind() == reflect.Ptr && fromField.IsNil() {
				continue
			}
			var sourceTime time.Time
			if fromField.Kind() == reflect.Ptr {
				sourceTime = fromField.Elem().Interface().(time.Time)
			} else {
				sourceTime = fromField.Interface().(time.Time)
			}

			if !pbField.CanSet() {
				continue
			}

			ts := timestamppb.New(sourceTime)
			pbField.Set(reflect.ValueOf(ts).Convert(pbField.Type()))

			continue
		}

		// Process sql.NullTime and *sql.NullTime fields
		if fromField.Type() == reflect.TypeOf(sql.NullTime{}) || fromField.Type() == reflect.TypeOf(&sql.NullTime{}) {
			if fromField.Kind() == reflect.Ptr && fromField.IsNil() {
				continue
			}
			var sourceNullTime sql.NullTime
			if fromField.Kind() == reflect.Ptr {
				sourceNullTime = fromField.Elem().Interface().(sql.NullTime)
			} else {
				sourceNullTime = fromField.Interface().(sql.NullTime)
			}

			if !pbField.CanSet() {
				continue
			}

			if sourceNullTime.Valid {
				ts := timestamppb.New(sourceNullTime.Time)
				pbField.Set(reflect.ValueOf(ts).Convert(pbField.Type()))
			} else {
				pbField.Set(reflect.Zero(pbField.Type()))
			}

			continue
		}
	}

	return nil
}

// getStructValue dereferences pointers and interfaces to obtain the underlying struct value.
// If obj is nil or not a struct, an invalid reflect.Value is returned.
func getStructValue(obj any) reflect.Value {
	value := reflect.ValueOf(obj)
	for value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}
