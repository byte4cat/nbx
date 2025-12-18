package dbu

import (
	"reflect"
	"testing"
)

func TestParseMDBUpdateData(t *testing.T) {
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
			"firstName": "John",
			"age":       30,
		}

		actual := BuildMongoUpdateMap(test, []string{"id", "lastName"})
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
