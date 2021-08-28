package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) map[string]string {
	errMap := map[string]string{}
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errMap[err.Field()] = fmt.Sprintf("%s is required", err.Field())
		case "gte":
			errMap[err.Field()] = fmt.Sprintf("%s value must be greater than %s", err.Field(), err.Param())
		}
	}
	return errMap

}
