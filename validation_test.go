package goutil

import "testing"

type User struct {
	Name      string `validate:"required"`
	Age       int    `validate:"min=10,max=40"`
	CreatedBy int    `validate:"required"`
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
		Age:       11,
		CreatedBy: 0,
	}

	if err := tValidation.ValidationStruct(user); err != nil {
		t.Log(err.Error())
	}
}
