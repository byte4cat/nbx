package dbu

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/byte4cat/nbx/v2/pkg/logger"
	"go.uber.org/zap"
)

// DefaultColumnNameFunc is the type for a function that provides a default column name
// for a struct field based on its Go name, used when no tags specify the name.
// Users can provide a custom function matching their GORM NamingStrategy.
type DefaultColumnNameFunc func(fieldName string) string

// DefaultSnakeCaseNamer is a DefaultColumnNameFunc that converts field names to snake_case.
// This serves as the built-in fallback if no custom NamingStrategy function is provided
// via SetDefaultColumnNameFunc or if SetDefaultColumnNameFunc(nil) is called.
func DefaultSnakeCaseNamer(fieldName string) string {
	return toSnakeCase(fieldName)
}

// packageDefaultNamer stores the package-level default column name function.
// It is initialized to DefaultSnakeCaseNamer.
var packageDefaultNamer DefaultColumnNameFunc = DefaultSnakeCaseNamer

// SetDefaultColumnNameFunc sets the package-level default column name function.
// This function should typically be called once during application initialization.
// If nil is passed, the default namer will be reset to DefaultSnakeCaseNamer.
// Note: If called concurrently with BuildRDBUpdateMapV6 after initialization,
// this could theoretically lead to race conditions, but this is an uncommon use case.
func SetDefaultColumnNameFunc(namer DefaultColumnNameFunc) {
	if namer == nil {
		packageDefaultNamer = DefaultSnakeCaseNamer
	} else {
		packageDefaultNamer = namer
	}
}

// BuildRDBUpdateMap constructs a map[string]any for use in relational database updates.
//
// This function recursively flattens the fields of a struct (or pointer to struct),
// including any fields marked with the GORM tag `embedded` (e.g., `gorm:"embedded;embeddedPrefix:foo_"`).
// For embedded fields, if an `embeddedPrefix` is specified, it is prepended to all keys from that embedded struct.
// The function supports both value and pointer embedded structs, and will skip nil pointers.
//
// The resulting map uses the following precedence for key names:
//  1. `gorm:"column:..."` tag
//  2. `bson:"..."` tag (if present and not "-")
//  3. For JSONB fields, `json:"..."` tag (if present and not "-")
//  4. Default naming strategy (snake_case by default)
//
// Fields listed in skipFields (with prefix applied) are omitted from the result.
// JSONB fields are marshaled to json.RawMessage.
//
// Returns an error if input is not a struct or pointer to struct, or if JSON marshaling fails.
func BuildRDBUpdateMap(x any, skipFields []string) (map[string]any, error) {
	result := make(map[string]any)
	skipMap := make(map[string]struct{}, len(skipFields))
	for _, field := range skipFields {
		skipMap[field] = struct{}{}
	}

	val := reflect.ValueOf(x)
	typ := reflect.TypeOf(x)

	// Handle nil input pointer
	if typ.Kind() == reflect.Pointer && val.IsNil() {
		return result, nil
	}

	// Dereference pointer if necessary
	if typ.Kind() == reflect.Pointer {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Ensure it's a struct
	if typ.Kind() != reflect.Struct {
		return result, fmt.Errorf("BuildRDBUpdateMapV6 input must be a struct or pointer to struct, got %s", typ.Kind())
	}

	var processingErr error

	// processFields recursively flattens struct fields into the result map.
	// If a field is marked with `gorm:"embedded"`, its fields are recursively processed.
	// If an embeddedPrefix is specified, it is prepended to all keys from that embedded struct.
	// Supports both value and pointer embedded structs.
	var processFields func(v reflect.Value, t reflect.Type, prefix string)
	processFields = func(v reflect.Value, t reflect.Type, prefix string) {
		if processingErr != nil {
			return
		}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldVal := v.Field(i)

			if !fieldVal.CanInterface() {
				continue
			}

			gormTagValue := field.Tag.Get("gorm")
			bsonTagValue := field.Tag.Get("bson")
			jsonTagValue := field.Tag.Get("json")

			logger.Debug("processing field", zap.String("fieldName", field.Name),
				zap.String("gormTag", gormTagValue), zap.String("bsonTag", bsonTagValue),
				zap.String("jsonTag", jsonTagValue),
			)

			// If the field is marked as embedded, recursively flatten its fields.
			// If an embeddedPrefix is specified, prepend it to all keys from this embedded struct.
			if strings.Contains(gormTagValue, "embedded") {
				logger.Debug("processing embedded", zap.String("fieldName", field.Name))
				embeddedPrefix := ""
				for part := range strings.SplitSeq(gormTagValue, ";") {
					part = strings.TrimSpace(part)
					if after, ok := strings.CutPrefix(part, "embeddedPrefix:"); ok {
						embeddedPrefix = after
						break
					}
				}
				if fieldVal.Kind() == reflect.Pointer {
					if fieldVal.IsNil() {
						continue
					}
					// Recursively process pointer embedded struct
					processFields(fieldVal.Elem(), fieldVal.Elem().Type(), prefix+embeddedPrefix)
				} else {
					// Recursively process value embedded struct
					processFields(fieldVal, field.Type, prefix+embeddedPrefix)
				}
				continue
			}

			var mapKey string
			isJSONB := strings.Contains(gormTagValue, "type:jsonb") || strings.Contains(gormTagValue, "jsonb")

			if strings.Contains(gormTagValue, "column:") {
				parts := strings.SplitSeq(gormTagValue, ";")
				for part := range parts {
					part = strings.TrimSpace(part)
					if after, ok := strings.CutPrefix(part, "column:"); ok {
						mapKey = after
						mapKey = strings.TrimSpace(mapKey)
						break
					}
				}
			}
			if mapKey == "" {
				if bsonTagValue != "" && bsonTagValue != "-" {
					mapKey = strings.Split(bsonTagValue, ",")[0]
				}
			}
			if mapKey == "" {
				if isJSONB {
					if jsonTagValue != "" && jsonTagValue != "-" {
						mapKey = strings.Split(jsonTagValue, ",")[0]
					}
				} else {
					mapKey = packageDefaultNamer(field.Name)
				}
			}
			if mapKey == "" {
				if isJSONB {
					if jsonTagValue != "" && jsonTagValue != "-" {
						mapKey = strings.Split(jsonTagValue, ",")[0]
					} else {
						mapKey = packageDefaultNamer(field.Name)
					}
				} else {
					mapKey = packageDefaultNamer(field.Name)
					if mapKey == "" && jsonTagValue != "" && jsonTagValue != "-" {
						mapKey = strings.Split(jsonTagValue, ",")[0]
					}
				}
			}
			if mapKey == "" {
				continue
			}
			if _, shouldSkip := skipMap[prefix+mapKey]; shouldSkip {
				continue
			}
			if isJSONB {
				if fieldVal.Kind() == reflect.Pointer && fieldVal.IsNil() {
					continue
				}
				valueToMarshal := fieldVal.Interface()
				if fieldVal.Kind() == reflect.Pointer {
					valueToMarshal = fieldVal.Elem().Interface()
				}
				jsonValue, err := json.Marshal(valueToMarshal)
				if err != nil {
					processingErr = fmt.Errorf("failed to marshal field %s (%s) to JSON: %w", field.Name, prefix+mapKey, err)
					return
				}
				result[prefix+mapKey] = json.RawMessage(jsonValue)
				continue
			}
			if fieldVal.Kind() == reflect.Pointer {
				if !fieldVal.IsNil() {
					result[prefix+mapKey] = fieldVal.Elem().Interface()
				}
			} else {
				result[prefix+mapKey] = fieldVal.Interface()
			}
		}
	}
	processFields(val, typ, "")
	return result, processingErr
}
