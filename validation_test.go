package goutil

import (
	"encoding/json"
	"testing"
)

type User struct {
	Name      string `json:"name" validate:"required"`
	Age       int    `json:"age" validate:"min=10,max=40"`
	CreatedBy int    `json:"created_by" validate:"required"`
}

var tValidation Validation

func init() {
	tValidation, _ = NewValidation()
}

// TestValidation how to run this process
// go test -v -run=TestValidation
func TestValidation(t *testing.T) {
	user := User{
		Name:      "Muhammad Rivaldy",
		Age:       1,
		CreatedBy: 0,
	}

	if validationErrors := tValidation.ValidationStruct(user); len(validationErrors.Errors) > 1 {
		errorInformations, _ := json.Marshal(validationErrors)
		t.Log(string(errorInformations))
	}
}
