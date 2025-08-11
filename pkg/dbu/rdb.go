package dbu

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/byte4cat/nbx/pkg/logger"
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
// It iterates over the fields of the input struct (or pointer to struct), extracts non-nil values,
// and builds a map suitable for database update operations.
// Unexported fields and fields with nil pointer values are skipped.
//
// The keys in the resulting map are determined by examining struct tags in a specific
// conditional priority order for each field:
//
// 1. gorm:"column:..." tag: The value specified in the 'column' option (e.g., `gorm:"column:my_col"`).
// 2. bson:"..." tag: The value before the first comma (e.g., `bson:"mongo_field,omitempty"`).
// 3. Conditional json:"..." tag (JSONB priority): If the field's gorm tag indicates it's a JSONB type (contains "type:jsonb" or "jsonb"), the json tag's value (e.g., `json:"api_field"`) is checked here. If a valid key is found, it's used. Otherwise (json tag missing or "-"), the process falls through to the default naming fallback.
// 4. Default Naming Fallback: If no previous tags (gorm:column, bson, and conditional json for JSONB) provided a valid key, the field name is converted using the package-level default naming function (see SetDefaultColumnNameFunc). By default, this is snake_case (see DefaultSnakeCaseNamer).
// 5. Final json:"..." tag (Non-JSONB fallback): If the field is NOT a JSONB type AND the default naming fallback (step 4) resulted in an empty key (which should only happen if the DefaultColumnNameFunc returns an empty string), the json tag's value is checked as a final option.
//
// Fields are skipped entirely if their derived map key is an empty string (e.g., a tag was set to "-")
// or if the derived key appears in the `skipFields` list.
//
// Anonymous embedded structs marked with `gorm:"embedded"` are recursively processed,
// and their internal fields are flattened into the top-level map according to the same
// key determination priority rules. This matches the standard GORM behavior for
// embedding in RDBs.
// Named embedded structs are treated as regular fields; the entire struct value will
// be included in the map if the field is not skipped (usually not suitable for RDB
// updates unless the target column is a compatible type like JSONB).
//
// Fields marked with `gorm:"type:jsonb"` or containing "jsonb" in their gorm tag
// are treated as JSONB columns. Their value is marshaled into a JSON string (`json.RawMessage`)
// before being added to the map.
//
// Parameters:
//
//	x: The struct or pointer to struct to traverse. Must be a struct or pointer to a non-nil struct.
//	skipFields: A list of strings representing map keys (database column names derived from tags or naming) to exclude from the result map.
//
// Returns:
//
//	A map[string]any containing the extracted fields and their values, intended for database updates.
//	An error if the input is not a valid struct or pointer, or if JSON marshaling of a JSONB field fails.
//
// Example:
//
//	type Address struct {
//		Street string `json:"street"` // json tag
//		City   string // No tags, relies on default namer
//	}
//	type User struct {
//		ID        uint   `gorm:"column:user_id"` // gorm column tag (highest priority)
//		FirstName string `json:"firstName"`      // json tag (lower prio for non-jsonb)
//		LastName  string // No tags, relies on default namer
//		address   Address `gorm:"embedded"`      // Anonymous embedded struct for flattening
//		Settings  map[string]string `gorm:"type:jsonb" json:"user_settings"` // JSONB + json tag
//	}
//
//	// Assume SetDefaultColumnNameFunc is set to DefaultSnakeCaseNamer in your init code.
//
//	user := User{
//		ID: 1,
//		FirstName: "Jane",
//		LastName: "Doe",
//		address: Address{Street: "Main St", City: "Anytown"},
//		Settings: map[string]string{"theme": "dark"},
//	}
//
//	updateMap, err := BuildRDBUpdateMap(user, []string{"user_id"})
//	// updateMap will be (assuming DefaultSnakeCaseNamer):
//	// {
//	//   "first_name": "Jane",            // Non-JSONB: defaultNamer("FirstName") -> "first_name" (wins over json:"firstName" for non-JSONB)
//	//   "last_name": "Doe",              // Non-JSONB: defaultNamer("LastName") -> "last_name"
//	//   "street": "Main St",             // Flattened, Non-JSONB: defaultNamer("Street") -> "street" (wins over json:"street")
//	//   "city": "Anytown",               // Flattened, Non-JSONB: defaultNamer("City") -> "city"
//	//   "user_settings": json.RawMessage(`{"theme":"dark"}`), // JSONB: json:"user_settings" -> "user_settings" (wins over defaultNamer for JSONB)
//	// }
//	// Note: "user_id" is skipped because it's in skipFields.
func BuildRDBUpdateMap(x any, skipFields []string) (map[string]any, error) {
	result := make(map[string]any)
	skipMap := make(map[string]struct{}, len(skipFields))
	for _, field := range skipFields {
		skipMap[field] = struct{}{}
	}

	val := reflect.ValueOf(x)
	typ := reflect.TypeOf(x)

	// Handle nil input pointer
	if typ.Kind() == reflect.Ptr && val.IsNil() {
		return result, nil // Return empty map for nil pointer, no error
	}

	// Dereference pointer if necessary
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// Ensure it's a struct
	if typ.Kind() != reflect.Struct {
		return result, fmt.Errorf("BuildRDBUpdateMapV6 input must be a struct or pointer to struct, got %s", typ.Kind())
	}

	var processingErr error // Use a variable to capture potential errors during recursion

	// Internal recursive function to process fields
	var processFields func(v reflect.Value, t reflect.Type)
	processFields = func(v reflect.Value, t reflect.Type) {
		// Stop processing if an error has already occurred
		if processingErr != nil {
			return
		}

		for i := range t.NumField() {
			field := t.Field(i)
			fieldVal := v.Field(i)

			// Skip unexported fields (cannot call .Interface() on them)
			if !fieldVal.CanInterface() {
				continue
			}

			// Get all relevant tag values
			gormTagValue := field.Tag.Get("gorm")
			bsonTagValue := field.Tag.Get("bson")
			jsonTagValue := field.Tag.Get("json")

			logger.Debug("processing field", zap.String("fieldName", field.Name),
				zap.String("gormTag", gormTagValue), zap.String("bsonTag", bsonTagValue),
				zap.String("jsonTag", jsonTagValue),
			)

			// Check for embedded struct (anonymous && has gorm:"embedded" tag)
			if field.Anonymous && strings.Contains(gormTagValue, "embedded") {
				logger.Debug("processing embedded", zap.String("fieldName", field.Name))
				// Recursively process embedded struct fields
				// Handle potential nil pointer embedded structs
				if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
					continue // Skip nil embedded pointers
				}
				if fieldVal.Kind() == reflect.Ptr {
					// Recursively process the element if it's a non-nil pointer
					processFields(fieldVal.Elem(), fieldVal.Elem().Type())
				} else {
					// Recursively process the value directly if non-pointer
					processFields(fieldVal, field.Type)
				}
				continue // Skip processing the embedded struct itself as a field
			}

			// --- Key Mapping Logic V6: Conditional Priority ---
			var mapKey string
			isJSONB := strings.Contains(gormTagValue, "type:jsonb") || strings.Contains(gormTagValue, "jsonb")

			// 1. Check gorm:"column:..." tag first (highest priority)
			if strings.Contains(gormTagValue, "column:") {
				parts := strings.SplitSeq(gormTagValue, ";")
				for part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "column:") {
						mapKey = strings.TrimPrefix(part, "column:")
						mapKey = strings.TrimSpace(mapKey)
						break
					}
				}
			}

			// 2. If gorm:"column" didn't provide a key, check bson tag
			if mapKey == "" {
				if bsonTagValue != "" && bsonTagValue != "-" {
					mapKey = strings.Split(bsonTagValue, ",")[0]
				}
			}

			// 3. Handle conditional json tag priority based on JSONB type
			if mapKey == "" {
				if isJSONB {
					// If JSONB, check json tag *before* default namer
					if jsonTagValue != "" && jsonTagValue != "-" {
						mapKey = strings.Split(jsonTagValue, ",")[0]
					}
					// If json tag is missing/empty for JSONB, will fall through to default namer below
				} else {
					// If NOT JSONB, apply default namer *before* checking json tag
					mapKey = packageDefaultNamer(field.Name)
					// If default namer returns empty, will fall through to json check below
				}
			}

			// 4. If still no key (means default namer returned empty OR it's non-JSONB), check json tag (lower priority for non-JSONB)
			if mapKey == "" {
				// This branch is reached if:
				// - gorm:column, bson failed AND (JSONB case: json failed)
				// - gorm:column, bson failed AND (NOT JSONB case: default namer failed)
				// Now check json tag as a final fallback for non-JSONB or if json was missing for JSONB
				// Re-check if it's JSONB because we handle it above.
				// This final json check should *only* apply if it's not JSONB OR if it is JSONB but the json tag was missing/empty
				// However, the logic `(if JSONB: json) -> defaultNamer -> (if NOT JSONB: json)` is awkward.
				// Let's refine step 3 & 4 logic.

				// Revised Step 3 & 4 Logic:
				// If mapKey is empty after gorm/bson:
				// If IsJSONB:
				//   Check json. If valid, use it.
				//   Else (json invalid/missing for JSONB): Use defaultNamer.
				// If !IsJSONB:
				//   Use defaultNamer.
				//   If defaultNamer returns empty: Check json. If valid, use it.

				// Let's rewrite the block if mapKey is empty after gorm/bson:
				if isJSONB {
					// JSONB case: json tag has higher priority than default namer
					if jsonTagValue != "" && jsonTagValue != "-" {
						mapKey = strings.Split(jsonTagValue, ",")[0]
					} else {
						// json tag missing/empty for JSONB field, fallback to default namer
						mapKey = packageDefaultNamer(field.Name)
					}
				} else {
					// NOT JSONB case: default namer has higher priority than json tag
					mapKey = packageDefaultNamer(field.Name)
					if mapKey == "" {
						// Default namer failed, check json tag as next fallback
						if jsonTagValue != "" && jsonTagValue != "-" {
							mapKey = strings.Split(jsonTagValue, ",")[0]
						}
					}
				}
				// End Revised Step 3 & 4 Logic
			}

			// If the map key is empty (e.g., a tag was "-"), skip it
			// Also skip if the default namer returned an empty string and no json fallback applied
			if mapKey == "" {
				continue
			}
			// --- End Key Mapping Logic V6 ---

			// Check if the field should be skipped based on the derived map key
			if _, shouldSkip := skipMap[mapKey]; shouldSkip {
				continue
			}

			// Handle JSONB fields based on gorm tag value
			// Note: This check happens *after* determining mapKey, using the full gormTagValue
			if isJSONB { // Use the pre-calculated isJSONB flag
				// Skip nil pointers for JSONB fields
				if fieldVal.Kind() == reflect.Ptr && fieldVal.IsNil() {
					continue
				}

				// Get the value to marshal (dereference if it's a pointer)
				valueToMarshal := fieldVal.Interface()
				if fieldVal.Kind() == reflect.Ptr {
					valueToMarshal = fieldVal.Elem().Interface()
				}

				// Marshal the value to JSON
				jsonValue, err := json.Marshal(valueToMarshal)
				if err != nil {
					// Capture the error and stop further processing in this path
					processingErr = fmt.Errorf("failed to marshal field %s (%s) to JSON: %w", field.Name, mapKey, err)
					return // Stop recursion on error
				}
				result[mapKey] = json.RawMessage(jsonValue)
				continue // Move to the next field after handling JSONB
			}

			// Handle standard fields (non-nil pointers and all non-pointers)
			if fieldVal.Kind() == reflect.Ptr {
				if !fieldVal.IsNil() {
					// Include non-nil pointers, store the dereferenced value
					result[mapKey] = fieldVal.Elem().Interface()
				}
				// Skip nil pointers explicitly
			} else {
				// Include non-pointer fields
				result[mapKey] = fieldVal.Interface()
			}
		}
	}

	// Start the processing from the top-level struct value
	processFields(val, typ)

	// Return the result map and any error encountered during processing
	return result, processingErr
}
