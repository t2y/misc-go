package main

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"

	"gopkg.in/go-playground/validator.v8"
)

type User struct {
	Name sql.NullString `binding:"required"`
	Age  sql.NullInt64  `binding:"required"`
}

type Users []*User

type Dept struct {
	Users Users `binding:"dive"`
}

var validate *validator.Validate

func main() {
	config := &validator.Config{TagName: "binding"}
	validate = validator.New(config)

	validate.RegisterCustomTypeFunc(
		ValidateValuer,
		sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})

	user := &User{
		Name: sql.NullString{String: "", Valid: true},
		Age:  sql.NullInt64{Int64: 0, Valid: false},
	}
	if err := validate.Struct(user); err != nil {
		fmt.Printf("Err:\n%+v\n", err)
	}

	users := make(Users, 0)
	users = append(users, user)
	users = append(users, user)
	fmt.Printf("users: %+v\n", users)

	for i := range users {
		if err := validate.Struct(users[i]); err != nil {
			fmt.Printf("Err:\n%+v\n", err)
		}
	}

	dept := &Dept{
		Users: users,
	}
	if err := validate.Struct(dept); err != nil {
		fmt.Printf("Err:\n%+v\n", err)
		fmt.Println("---")
		verrs := err.(validator.ValidationErrors)
		for key, ferr := range verrs {
			fmt.Printf("key: %s, field: %s, tag: %s\n", key, ferr.Field, ferr.Tag)
		}
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
