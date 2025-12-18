package dbu

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRDBUpdateData(t *testing.T) {
	// set default naming strategy
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

	t.Run("[SUCCESS] should skip fields and returns the correct values", func(t *testing.T) {
		type testStruct struct {
			ID        string `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Age       int    `json:"age"`
		}

		test := testStruct{
			ID:        "123",
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
		}

		expected := map[string]any{
			"first_name": "John",
			"age":        30,
		}

		actual, err := BuildRDBUpdateMap(test, []string{"id", "last_name"})
		require.NoError(t, err)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("[SUCCESS] should return map[string]any with the correct values", func(t *testing.T) {
		type testStruct struct {
			ID        string `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Age       int    `json:"age"`
		}

		test := testStruct{
			ID:        "123",
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
		}

		expected := map[string]any{
			"id":         "123",
			"first_name": "John",
			"last_name":  "Doe",
			"age":        30,
		}

		actual, err := BuildRDBUpdateMap(test, []string{})

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("[SUCCESS] should return map[string]any with the correct values when some value is pointer", func(t *testing.T) {
		type testStruct struct {
			ID        *string `json:"id"`
			FirstName *string `json:"firstName"`
			LastName  string  `json:"lastName"`
			Age       int     `json:"age"`
		}

		id := "123"
		firstName := "John"
		test := testStruct{
			ID:        &id,
			FirstName: &firstName,
			LastName:  "Doe",
			Age:       30,
		}

		expected := map[string]any{
			"id":         "123",
			"first_name": "John",
			"last_name":  "Doe",
			"age":        30,
		}

		actual, err := BuildRDBUpdateMap(test, []string{})
		require.NoError(t, err)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("[SUCCESS] should return map[string]any with the correct values when some value is pointer", func(t *testing.T) {
		type testStruct struct {
			ID        *string `json:"id"`
			FirstName *string `json:"firstName"`
			LastName  *string `json:"lastName"`
			Age       *int    `json:"age"`
		}

		id := "123"
		firstName := "John"
		lastName := "Doe"
		age := 30
		test := testStruct{
			ID:        &id,
			FirstName: &firstName,
			LastName:  &lastName,
			Age:       &age,
		}

		expected := map[string]any{
			"id":         "123",
			"first_name": "John",
			"last_name":  "Doe",
			"age":        30,
		}

		actual, err := BuildRDBUpdateMap(test, []string{})
		require.NoError(t, err)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("[SUCCESS] should return map[string]any with the correct values when embedded is set", func(t *testing.T) {
		type Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}

		type user struct {
			ID        string `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Age       int    `json:"age"`
			Address   `gorm:"embedded" json:"address"`
		}

		test := user{
			ID:        "123",
			FirstName: "John",
			LastName:  "Doe",
			Age:       30,
			Address: Address{
				Street: "street",
				City:   "city",
			},
		}

		expected := map[string]any{
			"id":         "123",
			"first_name": "John",
			"last_name":  "Doe",
			"age":        30,
			"street":     "street",
			"city":       "city",
		}

		actual, err := BuildRDBUpdateMap(test, []string{})
		require.NoError(t, err)

		// t.Logf("expected: %v", expected)
		// t.Logf("actual: %v", actual)

		assert.Equal(t, expected, actual)

	})

	t.Run("[SUCCESS] should return map[string]any with the correct values when embedded is set and some value is pointer", func(t *testing.T) {
		type Address struct {
			Street  *string `json:"street"`
			City    string  `json:"city"`
			ZipCode *string `json:"zipCode"`
		}

		type testStruct struct {
			ID        string  `json:"id"`
			FirstName *string `json:"firstName"`
			LastName  string  `json:"lastName"`
			Age       *int    `json:"age"`
			Address   `gorm:"embedded" json:"address"`
		}

		name := "John"
		street := "street"
		test := testStruct{
			ID:        "123",
			FirstName: &name,
			LastName:  "Doe",
			Age:       nil,
			Address: Address{
				Street:  &street,
				City:    "city",
				ZipCode: nil,
			},
		}

		expected := map[string]any{
			"id":         "123",
			"first_name": "John",
			"last_name":  "Doe",
			"street":     "street",
			"city":       "city",
		}

		actual, err := BuildRDBUpdateMap(test, []string{})
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})

	t.Run("[SUCCESS] should return map[string]any with the correct values when jsonb is set and some value is pointer", func(t *testing.T) {
		type address struct {
			Street  *string `json:"street"`
			City    string  `json:"city"`
			ZipCode *string `json:"zipCode"`
		}

		type testStruct struct {
			ID          string  `json:"id"`
			FirstName   *string `json:"firstName"`
			LastName    string  `json:"lastName"`
			Age         *int    `json:"age"`
			Origin      address `gorm:"type:jsonb" json:"origin"`
			Destination address `gorm:"type:jsonb" json:"destination"`
		}

		street := "123 Test Street"
		zip := "90001"
		firstName := "John"
		age := 30

		input := testStruct{
			ID:        "abc-123",
			FirstName: &firstName,
			LastName:  "Doe",
			Age:       &age,
			Origin: address{
				Street:  &street,
				City:    "Los Angeles",
				ZipCode: &zip,
			},
			Destination: address{
				Street:  nil,
				City:    "San Francisco",
				ZipCode: nil,
			},
		}

		expected := map[string]any{
			"id":         "abc-123",
			"first_name": "John",
			"last_name":  "Doe",
			"age":        30,
			"origin": json.RawMessage(`{
				"street": "123 Test Street",
				"city": "Los Angeles",
				"zipCode": "90001"
			}`),
			"destination": json.RawMessage(`{
				"street": null,
				"city": "San Francisco",
				"zipCode": null
			}`),
		}

		actual, err := BuildRDBUpdateMap(input, []string{})
		require.NoError(t, err)

		expJSON, _ := json.Marshal(decodeRawMessages(expected))
		actJSON, _ := json.Marshal(decodeRawMessages(actual))

		assert.Equal(t, string(expJSON), string(actJSON))
	})
}

func Test_toSnakeCase(t *testing.T) {
	t.Run("[SUCCESS] should return snake case string", func(t *testing.T) {

		actual := []string{
			toSnakeCase(""),
			toSnakeCase("Name"),
			toSnakeCase("Role"),
			toSnakeCase("role"),
			toSnakeCase("IsEnable"),
			toSnakeCase("IsDisable"),
			toSnakeCase("isDisable"),
			toSnakeCase("IsDeleted"),
			toSnakeCase("isDeleted"),
			toSnakeCase("CreatedAt"),
			toSnakeCase("UpdatedAt"),
			toSnakeCase("DeletedAt"),
			toSnakeCase("deletedAt"),
			toSnakeCase("IsDeleted"),
			toSnakeCase("isDeleted"),
			toSnakeCase("UserID"),
			toSnakeCase("FirstName"),
			toSnakeCase("FirstNameLastName"),
			toSnakeCase("FirstNameLastNameAge"),
			toSnakeCase("FirstNameLastNameAgeID"),
			toSnakeCase("FirstNameLastNameAgeEmployeeID"),
		}

		expected := []string{
			"",
			"name",
			"role",
			"role",
			"is_enable",
			"is_disable",
			"is_disable",
			"is_deleted",
			"is_deleted",
			"created_at",
			"updated_at",
			"deleted_at",
			"deleted_at",
			"is_deleted",
			"is_deleted",
			"user_id",
			"first_name",
			"first_name_last_name",
			"first_name_last_name_age",
			"first_name_last_name_age_id",
			"first_name_last_name_age_employee_id",
		}

		for i, a := range actual {
			if a != expected[i] {
				assert.Equal(t, expected[i], a)
			}
		}
	})

	t.Run("[SUCCESS] should correctly parse nationalIDNo to national_id_no", func(t *testing.T) {
		type testStruct struct {
			NationalIDNo string `json:"nationalIDNo"`
		}

		input := testStruct{
			NationalIDNo: "A123456789",
		}

		expected := map[string]any{
			"national_id_no": "A123456789",
		}

		actual, err := BuildRDBUpdateMap(input, []string{})
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}

func TestBuildRDBUpdateMap_EmbeddedPrefix(t *testing.T) {
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

	type Details struct {
		Text string `json:"text"`
		PDF  string `json:"pdf"`
	}

	type Travel struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Details `gorm:"embedded;embeddedPrefix:details_" json:"details"`
	}

	obj := Travel{
		ID:   "123",
		Name: "Taipei Trip",
		Details: Details{
			Text: "hello",
			PDF:  "url.pdf",
		},
	}

	expected := map[string]any{
		"id":           "123",
		"name":         "Taipei Trip",
		"details_text": "hello",
		"details_pdf":  "url.pdf",
	}

	result, err := BuildRDBUpdateMap(obj, nil)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBuildRDBUpdateMap_EmbeddedPrefixWithPointer(t *testing.T) {
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

	type Details struct {
		Text *string `json:"text"`
		PDF  *string `json:"pdf"`
	}

	type Travel struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Details *Details `gorm:"embedded;embeddedPrefix:details_" json:"details"`
	}

	text := "the text"
	pdf := "file.pdf"
	obj := Travel{
		ID:   "999",
		Name: "Kaohsiung",
		Details: &Details{
			Text: &text,
			PDF:  &pdf,
		},
	}

	expected := map[string]any{
		"id":           "999",
		"name":         "Kaohsiung",
		"details_text": "the text",
		"details_pdf":  "file.pdf",
	}

	result, err := BuildRDBUpdateMap(obj, nil)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBuildRDBUpdateMap_EmbeddedPrefix_SkipNilPointer(t *testing.T) {
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

	type Details struct {
		Text *string `json:"text"`
		PDF  *string `json:"pdf"`
	}

	type Travel struct {
		ID      string   `json:"id"`
		Name    string   `json:"name"`
		Details *Details `gorm:"embedded;embeddedPrefix:details_" json:"details"`
	}

	// Details 為 nil，預期 details_ 欄位都不會出現
	obj := Travel{
		ID:      "111",
		Name:    "No Details",
		Details: nil,
	}

	expected := map[string]any{
		"id":   "111",
		"name": "No Details",
	}

	result, err := BuildRDBUpdateMap(obj, nil)
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBuildRDBUpdateMap_EmbeddedPrefix_ShouldPanicIfDefective(t *testing.T) {
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

	// 1. Simulate an embedded struct (UpdateTravelDetailRequest)
	type Details struct {
		Text *string `json:"text" validate:"omitempty,max=5000"`
		PDF  *string `json:"pdf" validate:"omitempty"`
	}

	// 2. Simulate the parent struct (UpdateTravelRequest)
	type Travel struct {
		ID       string  `json:"id"`
		Name     *string `json:"name"` // Contains other pointer fields
		IsPublic *bool   `json:"isPublic"`
		// 🚨 Critical field: embedded pointer
		Details *Details `gorm:"embedded;embeddedPrefix:details_" json:"details"`
	}

	t.Run("[FAILURE_SIM] Should cause panic if DBU library is defective when Details is nil", func(t *testing.T) {

		name := "Taipei Trip New"

		obj := Travel{
			ID:       "123",
			Name:     &name,
			IsPublic: nil,
			Details:  nil, // Simulate update request where Details is nil
		}

		expected := map[string]any{
			"id":   "123",
			"name": "Taipei Trip New",
			// Expected: no fields from Details should appear
		}

		// Attempt to capture panic. If the DBU library is defective, this should crash.
		// NOTE: In Go tests, if a panic occurs inside the test function, the test fails.
		// We use defer + recover here to detect a panic without aborting the test prematurely.

		var actual map[string]any
		var err error

		// Since we want to simulate a panic, we wrap the call in an anonymous function
		// and use recover() to safely detect it so we can proceed to assertions.
		didPanic := func() (panicked bool) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Successfully reproduced panic: %v", r)
					panicked = true
				}
			}()

			actual, err = BuildRDBUpdateMap(obj, nil)
			return false
		}()

		if didPanic {
			// If a panic occurred, it means the DBU library indeed has a bug
			// when handling nil embedded pointers.
			t.Errorf("DBU library is defective: BuildRDBUpdateMap panicked when Details was nil")
			t.FailNow()
		}

		// If no panic occurred, proceed to verify correctness (nil fields should be skipped)
		if err != nil {
			t.Fatalf("BuildRDBUpdateMap returned an unexpected error: %v", err)
		}

		// Verify the result matches expectations (Details fields should be omitted)
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}

		t.Log("result", actual)
	})
}

func TestBuildRDBUpdateMap_ShouldSkipNilTimeField(t *testing.T) {
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

	type Details struct {
		Text *string `json:"text"`
		PDF  *string `json:"pdf"`
	}

	type Travel struct {
		ID          string     `json:"id"`
		Name        *string    `json:"name"`
		IsPublic    *bool      `json:"isPublic"`
		StartAt     *time.Time `json:"startAt"`
		EndAt       *time.Time `json:"endAt"`
		Details     *Details   `gorm:"embedded;embeddedPrefix:details_" json:"details"`
		Category    int        `json:"category"`
		Destination string     `json:"destination"`
	}

	// Simulate input with StartAt = nil
	name := "Test Trip"
	endAt := time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC)

	obj := Travel{
		ID:          "test123",
		Name:        &name,
		IsPublic:    nil,
		StartAt:     nil, // the field we want to test
		EndAt:       &endAt,
		Details:     nil,
		Category:    1,
		Destination: "Taipei",
	}

	expected := map[string]any{
		"id":          "test123",
		"name":        "Test Trip",
		"end_at":      endAt,
		"category":    1,
		"destination": "Taipei",
		// StartAt should be skipped because it is nil
	}

	actual, err := BuildRDBUpdateMap(obj, nil)
	if err != nil {
		t.Fatalf("BuildRDBUpdateMap returned unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected update map: %+v, got: %+v", expected, actual)
	}

	t.Logf("DBU produced update map: %+v", actual)
}

func decodeRawMessages(m map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range m {
		if raw, ok := v.(json.RawMessage); ok {
			var decoded any
			_ = json.Unmarshal(raw, &decoded)
			out[k] = decoded
		} else {
			out[k] = v
		}
	}
	return out
}
