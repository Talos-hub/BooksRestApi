package validations

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
)

const minimallyFields = 5

const (
	fieldId     = "id"
	fieldTitle  = "title"
	fieldGenre  = "genre"
	fieldAuthor = "author"
)

var (
	sqlInjectionRegex = regexp.MustCompile(`(?i)(\b(UNION|SELECT|INSERT|DELETE|UPDATE|DROP|ALTER|CREATE|EXEC)\b|--|;|/\*|\*/|xp_)`)
	xssRegex          = regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onload=|onclick=)`)
)

// Warning
// You could say why I don't use simple way like package validator or type Validator interface{Valudate()error}
// But I want to to know how works with reflection proparly and safety.
// It shows fundamentaion data about golang, and proparly aproach with refelction
// As well as I don't use meta tegs because it does't conveniant in this case.
// Also it has very good documentation and it doesn't needs any dependencies.

// Validate might validates structions which contains fileds like:
// ID, Title, Ganre, Author, PublicationData.
// If a type don't conatins all these fields or the type is not a sturct
// It returns message error like: expected a struct
// Also there is reflection, it might be slow.
// It is safty function, if you will try to use parameter like a function and stuff It don't panic
func Validate(book any) error {
	// it cannot validate nil object
	// book is nil when type and value of inteface{} are nil
	if book == nil {
		return apperrors.NewValidateErr("canot validate nil object", nil,
			apperrors.NewValidateReflectErr("error validation, type is nil"))
	}

	// just get value of this interface
	// describe
	// interface{} has two things: type and value
	// describe
	value := reflect.ValueOf(book)

	// Handle pointers by dereferencing them
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return apperrors.NewValidateErr("cannot validate nil pointer", nil,
				apperrors.NewValidateReflectErr("error validation, type is nil pointer"))
		}
		// dereferencing
		value = value.Elem()
	}

	// check that interface is struction
	// for instance it might be func() or somthing else
	// So it needs to check that interface is struction
	if value.Kind() != reflect.Struct {
		return apperrors.NewValidateErr("cannot validate non-struction type", nil,
			apperrors.NewValidateReflectErr("error validation, type is not struction"))
	}

	// it use minimallyFileds for avoid allocation
	validationErrors := make([]string, 0, minimallyFields)

	validationErrors = append(validationErrors, validateBookFields(value)...)

	if len(validationErrors) > 0 {
		return apperrors.NewValidateErr("error validation", validationErrors, errors.New("error validation"))
	}
	// TODO
	return nil
}

// validateBookFields checks fields and chooses
// a specific function for validation
// if spicific fields is wrong it add it to slice that contains errors message
func validateBookFields(value reflect.Value) []string {
	// here is errors message from validation functions
	// when validation is finished it function returns the slice to up
	errorsSlice := make([]string, 0, value.NumField())
	if value.NumField() == 0 {
		errorsSlice = append(errorsSlice, "struct is empty")
		return errorsSlice
	}

	t := value.Type()

	for i := 0; i < value.NumField(); i++ {
		// get type of a field
		field := t.Field(i)

		// avoid panic, because if a field is not exported it can be panic
		if !field.IsExported() {
			continue
		}

		// get value of a field
		fieldValue := value.Field(i)

		//handle pointers by dereferencing them
		if fieldValue.Kind() == reflect.Pointer {
			// if a pointer is nil add to a error slice
			if fieldValue.IsNil() {
				errorsSlice = append(errorsSlice, fmt.Sprintf("%s field is nil", field.Name))
				continue
			}
			// dereferencing
			fieldValue = fieldValue.Elem()
		}

		// handle nested structions recursively (skip time.Time since the time is valid)
		if fieldValue.Kind() == reflect.Struct {
			// time are valid
			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				continue
			}
			// recursively
			nestedErrs := validateBookFields(fieldValue)
			errorsSlice = append(errorsSlice, nestedErrs...)
			continue
		}

		// Validate based on field name (case-insensitive)
		fieldName := strings.ToLower(field.Name)

		switch fieldName {
		case fieldId:
			errorsSlice = append(errorsSlice, validateID(fieldValue)...)
		case fieldTitle:
			errorsSlice = append(errorsSlice, validateString(fieldValue, fieldName)...)
		case fieldGenre:
			errorsSlice = append(errorsSlice, validateString(fieldValue, fieldName)...)
		case fieldAuthor:
			errorsSlice = append(errorsSlice, validateString(fieldValue, fieldName)...)
		}
	}
	return errorsSlice

}

// field specific validation function
func validateID(value reflect.Value) []string {
	if value.Kind() != reflect.Uint64 {
		return []string{"id must be uint64"}
	}
	if value.Uint() == 0 {
		return []string{"id cannot be ziro"}
	}
	return nil
}

// validateString is specific function that validate string type.
// When nameField: Author or Gangre etc
func validateString(value reflect.Value, nameField string) []string {
	if value.Kind() != reflect.String {
		return []string{fmt.Sprintf("%s: must be string", nameField)}
	}
	if value.String() == "" {
		return []string{fmt.Sprintf("%s: cannot be empty", nameField)}
	}
	if len(value.String()) > 100 {
		return []string{fmt.Sprintf("%s: cannot be large than 100", nameField)}
	}
	// check malicious injections
	message := validateMaliciousInjections(value.String())
	if len(message) > 0 {
		return []string{fmt.Sprintf("field: %s, %s", nameField, message)}
	}
	return nil
}

// validateMaliciousInjection if that instance is "clean",
// it returns empty string, if not it returns message with info about
// malicious pattern
func validateMaliciousInjections(str string) string {
	if sqlInjectionRegex.MatchString(str) {
		return fmt.Sprintf("%s contains SQL injection patterns", str)
	}
	if xssRegex.MatchString(str) {
		return fmt.Sprintf("%s contatins XsS pattern", str)
	}

	return ""

}
