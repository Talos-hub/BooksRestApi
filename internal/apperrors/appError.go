package apperrors

import (
	"fmt"
	"strings"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {

		return fmt.Sprintf("%s, %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
func NewAppError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

type ValidateErr struct {
	Message string
	Fields  []string
	Err     error
}

func (v *ValidateErr) Error() string {
	if v == nil {
		return ""
	}
	var sb strings.Builder

	// Use the Message field if provided
	if v.Message != "" {
		sb.WriteString(v.Message)
		sb.WriteString("\n")
	}

	// Add field errors with proper formatting
	if len(v.Fields) > 0 {
		sb.WriteString("Validation errors:\n")
		for _, field := range v.Fields {
			sb.WriteString(fmt.Sprintf("  • %s\n", field))
		}
	}

	// Add underlying error with separation
	if v.Err != nil {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("Underlying error: %s", v.Err.Error()))
	}
	return sb.String()

}

func (v *ValidateErr) Unwrap() error {
	return v.Err
}

func NewValidateErr(mesage string, fields []string, err error) *ValidateErr {
	return &ValidateErr{
		Message: mesage,
		Fields:  fields,
		Err:     err,
	}
}

type ValidationReflectErr struct {
	Message string
}

func (v *ValidationReflectErr) Error() string {
	return v.Message
}

func NewValidateReflectErr(message string) *ValidateErr {
	return &ValidateErr{Message: message}
}
