package dbu

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yimincai/nbx/pkg/logger"
)

func TestMain(m *testing.M) {
	// setup
	_, err := logger.New(logger.DefaultConfig())
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}

	os.Exit(m.Run())
}

func TestParseMDBUpdateData(t *testing.T) {
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
			"id":        "123",
			"firstName": "John",
			"lastName":  "Doe",
			"age":       30,
		}

		actual := BuildMongoUpdateMap(test, []string{})
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
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
			"id":        "123",
			"firstName": "John",
			"lastName":  "Doe",
			"age":       30,
		}

		actual := BuildMongoUpdateMap(test, []string{})
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

func TestParseRDBUpdateData(t *testing.T) {

	// set default naming strategy
	SetDefaultColumnNameFunc(DefaultSnakeCaseNamer)

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
			Address   `gorm:"embedded" jjson:"address"`
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

func BenchmarkParseMDBUpdateData(b *testing.B) {
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

	for range 1000000 {
		BuildMongoUpdateMap(test, []string{})
	}
}

func BenchmarkParseRDBUpdateData(b *testing.B) {
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

	for range 1000000 {
		BuildRDBUpdateMap(test, []string{})
	}
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
