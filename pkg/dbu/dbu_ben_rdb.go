package dbu

import "testing"

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
