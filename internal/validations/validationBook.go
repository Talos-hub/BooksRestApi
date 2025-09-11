package validations

import (
	"reflect"

	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
)

const minimallyFields = 5

// Warning
// You could say why I don't use simple way like package validator or type Validator interface{Valudate()error}
// But I want to to know how works with reflection proparly and safety.
// It shows fundamentaion data about golang, and proparl aproach with refelction
// As well as I don't use meta tegs because it does't conveniant in this case.
// Also it has very good documentation and it doesn't needs any dependencies.

// Validate might validates structions which contains fileds like:
// ID, Title, Ganre, Author, PublicationData.
// If a type don't conatins all these fields or the type is not a sturct
// It returns message error like: expected a struct
// Also there is reflection, it might be slow.
// It is safty function, if you will try to use parameter like a function and stuff It don't panic
func Validate(book any) *apperrors.ValidateErr {
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
	// TODO validationErrors := make([]string, 0, minimallyFields)

	// TODO
	return nil
}
