package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v8"
)

type User struct {
	Name sql.NullString `validate:"required"`
	Age  sql.NullInt64  `validate:"required"`
}

type Users []*User

var validate *validator.Validate

func main() {
	config := &validator.Config{TagName: "validate"}
	validate = validator.New(config)

	validate.RegisterCustomTypeFunc(
		ValidateValuer,
		sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})

	x := User{
		Name: sql.NullString{String: "", Valid: true},
		Age:  sql.NullInt64{Int64: 0, Valid: false},
	}
	if err := validate.Struct(x); err != nil {
		fmt.Printf("Err:\n%+v\n", err)
	}
}

func ValidateValuer(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}

	return nil
}
