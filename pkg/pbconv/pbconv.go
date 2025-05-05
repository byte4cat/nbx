package pbconv

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type pathNode struct {
	FieldName string
	SubPath   []*pathNode
}

var defaultFieldsMap = map[string]struct{}{
	"CreatedAt": {},
	"UpdatedAt": {},
	"DeletedAt": {},
	"StartAt":   {},
	"EndAt":     {},
	"Start":     {},
	"End":       {},
}

var defaultFields = []string{"CreatedAt", "UpdatedAt", "DeletedAt", "StartAt", "EndAt", "Start", "End"}

// StructTimeToPbTimestamp converts specified time.Time fields in fromObj
// to *timestamppb.Timestamp fields in pbObj. If no field paths are provided,
// a set of default field names is used.
func StructTimeToPbTimestamp(pbObj, fromObj any, fieldPaths ...string) error {
	if len(fieldPaths) == 0 {
		fieldPaths = defaultFields
	}
	return processStruct(getStructValue(pbObj), getStructValue(fromObj), parseFieldPaths(fieldPaths))
}

// SliceStructTimeToPbTimestamp converts time.Time fields in a slice of structs (fromObjSlice)
// to *timestamppb.Timestamp fields in the corresponding protobuf object slice (pbObjSlice).
// It supports nested field paths such as "Event.StartAt" or "Event.CreatedAt".
func SliceStructTimeToPbTimestamp(pbObjSlice any, fromObjSlice any, fieldPaths ...string) error {
	pbObjType := reflect.TypeOf(pbObjSlice)
	fromObjType := reflect.TypeOf(fromObjSlice)
	if pbObjType.Kind() != reflect.Ptr || pbObjType.Elem().Kind() != reflect.Slice ||
		fromObjType.Kind() != reflect.Ptr || fromObjType.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("pbObjSlice and fromObjSlice must be slices")
	}

	pbObjSliceValue := reflect.ValueOf(pbObjSlice).Elem()
	fromObjSliceValue := reflect.ValueOf(fromObjSlice).Elem()

	if pbObjSliceValue.Len() != fromObjSliceValue.Len() {
		return fmt.Errorf("pbObjSlice and fromObjSlice must be the same length")
	}

	for i := range pbObjSliceValue.Len() {
		pbObj := pbObjSliceValue.Index(i).Addr().Interface()
		fromObj := fromObjSliceValue.Index(i).Addr().Interface()

		err := StructTimeToPbTimestamp(pbObj, fromObj, fieldPaths...)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseFieldPaths(paths []string) []*pathNode {
	tree := map[string]*pathNode{}
	for _, path := range paths {
		parts := strings.Split(path, ".")
		buildPathTree(tree, parts)
	}
	return mapValues(tree)
}

func buildPathTree(tree map[string]*pathNode, parts []string) {
	if len(parts) == 0 {
		return
	}
	head := parts[0]
	node, exists := tree[head]
	if !exists {
		node = &pathNode{FieldName: head}
		tree[head] = node
	}
	buildPathTree(mapify(node.SubPath), parts[1:])
	node.SubPath = mapValues(mapify(node.SubPath))
}

func mapify(nodes []*pathNode) map[string]*pathNode {
	m := make(map[string]*pathNode)
	for _, n := range nodes {
		m[n.FieldName] = n
	}
	return m
}

func mapValues(m map[string]*pathNode) []*pathNode {
	values := make([]*pathNode, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func processStruct(pbValue, fromValue reflect.Value, paths []*pathNode) error {
	if !pbValue.IsValid() || !fromValue.IsValid() {
		return nil
	}

	for _, node := range paths {
		pbf := pbValue.FieldByName(node.FieldName)
		ff := fromValue.FieldByName(node.FieldName)

		if !pbf.IsValid() || !ff.IsValid() {
			continue
		}

		if len(node.SubPath) > 0 {
			if pbf.Kind() == reflect.Slice && ff.Kind() == reflect.Slice {
				for j := range pbf.Len() {
					err := processStruct(
						getStructValue(pbf.Index(j)),
						getStructValue(ff.Index(j)),
						node.SubPath,
					)
					if err != nil {
						return fmt.Errorf("%s[%d]: %w", node.FieldName, j, err)
					}
				}
			} else {
				err := processStruct(getStructValue(pbf), getStructValue(ff), node.SubPath)
				if err != nil {
					return fmt.Errorf("%s: %w", node.FieldName, err)
				}
			}
			continue
		}

		// no subpath, check if the field is a time.Time or *time.Time
		if _, ok := defaultFieldsMap[node.FieldName]; !ok {
			continue
		}
		if err := convertTimeField(pbf, ff); err != nil {
			return fmt.Errorf("field %s: %w", node.FieldName, err)
		}
	}
	return nil
}

func convertTimeField(pbField, fromField reflect.Value) error {
	if !fromField.IsValid() || !pbField.IsValid() {
		return nil
	}

	isTime := fromField.Type() == reflect.TypeOf(time.Time{}) || fromField.Type() == reflect.TypeOf(&time.Time{})
	if !isTime {
		return nil
	}

	if fromField.Kind() == reflect.Ptr && fromField.IsNil() {
		return nil
	}

	var t time.Time
	if fromField.Kind() == reflect.Ptr {
		t = fromField.Elem().Interface().(time.Time)
	} else {
		t = fromField.Interface().(time.Time)
	}

	ts := timestamppb.New(t)
	if pbField.CanSet() {
		pbField.Set(reflect.ValueOf(ts).Convert(pbField.Type()))
	}
	return nil
}

func getStructValue(obj any) reflect.Value {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return reflect.Value{}
		}
		v = v.Elem()
	}
	return v
}
