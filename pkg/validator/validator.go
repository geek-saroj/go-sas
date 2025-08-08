package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Use JSON tags instead of struct field names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates any struct and returns field-specific error messages
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Collect validation errors
	errors := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field() // will now be json tag
			tag := e.Tag()

			switch tag {
			case "required":
				errors[field] = "This field is required"
			case "email":
				errors[field] = "Invalid email format"
			case "min":
				errors[field] = fmt.Sprintf("Minimum length is %s", e.Param())
			case "max":
				errors[field] = fmt.Sprintf("Maximum length is %s", e.Param())
			case "len":
				errors[field] = fmt.Sprintf("Length must be %s", e.Param())
			default:
				errors[field] = fmt.Sprintf("Invalid value for %s", field)
			}
		}
	} else {
		errors["error"] = "Internal validation error"
	}
	return errors
}
