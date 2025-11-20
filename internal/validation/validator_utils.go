package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func Errors(fieldErrors validator.ValidationErrors) []string {
	var displayErrors []string
	for _, fieldError := range fieldErrors {
		if fieldError.Tag() == "required" {
			displayErrors = append(displayErrors, fmt.Sprintf("%v is %v", fieldError.Namespace(), fieldError.Tag()))
		}
	}
	return displayErrors
}

func ErrorsAsString(fieldErrors validator.ValidationErrors) string {
	return strings.Join(Errors(fieldErrors), ",")
}
