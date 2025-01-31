package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
	"strings"
)

func ExtractValidationErrors(req interface{}) (err error) {
	var builder strings.Builder
	var v = NewValidator()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err = v.Struct(req)
	if err != nil {
		for i, err := range err.(validator.ValidationErrors) {
			if i > 0 {
				builder.WriteString(" | ")
			}
			builder.WriteString(err.Field() + ": " + err.Tag())
		}

		return errors.New(builder.String())
	}

	return nil
}

func RequiredIf(fl validator.FieldLevel) bool {
	childField := fl.Field()
	params := strings.Split(fl.Param(), ":")
	if len(params) != 2 {
		return false
	}

	fieldName, condition := params[0], params[1]

	parentField := fl.Parent().FieldByName(fieldName)
	if !childField.IsValid() && parentField.String() == condition {
		return !childField.IsZero()
	}

	return true
}

func NewValidator() *validator.Validate {
	v := validator.New()

	err := v.RegisterValidation("required_if", RequiredIf)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to register 'required_if' validation: %v", err))
	}
	return v
}
