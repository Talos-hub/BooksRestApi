package validations

import (
	"reflect"
	"testing"
	"time"

	"github.com/Talos-hub/BooksRestApi/internal/apperrors"
)

// Test structures
type ValidBook struct {
	ID              uint64
	Title           string
	Genre           string
	Author          string
	PublicationDate time.Time
}

type BookWithPointers struct {
	ID              *uint64
	Title           *string
	Genre           *string
	Author          *string
	PublicationDate *time.Time
}

type BookWithUnexported struct {
	ID              uint64
	Title           string
	Genre           string
	Author          string
	publicationDate time.Time // unexported field
}

type NestedBook struct {
	Book      ValidBook
	Metadata  map[string]string
	CreatedAt time.Time
}

type InvalidStruct struct {
	ID    string // wrong type
	Title int    // wrong type
}

// Test data
var (
	validUint64 = uint64(1)
	validString = "Valid String"
	validTime   = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	invalidTime = time.Time{}
)

func TestValidate_ValidBook(t *testing.T) {
	book := ValidBook{
		ID:              1,
		Title:           "Clean Code",
		Genre:           "Programming",
		Author:          "Robert C. Martin",
		PublicationDate: validTime,
	}

	err := Validate(book)
	if err != nil {
		t.Errorf("Expected nil error for valid book, got: %v", err)
	}
}

func TestValidate_ValidBookPointer(t *testing.T) {
	book := &ValidBook{
		ID:              1,
		Title:           "Clean Code",
		Genre:           "Programming",
		Author:          "Robert C. Martin",
		PublicationDate: validTime,
	}

	err := Validate(book)
	if err != nil {
		t.Errorf("Expected nil error for valid book pointer, got: %v", err)
	}
}

func TestValidate_ValidBookWithPointers(t *testing.T) {
	book := BookWithPointers{
		ID:              &validUint64,
		Title:           &validString,
		Genre:           &validString,
		Author:          &validString,
		PublicationDate: &validTime,
	}

	err := Validate(book)
	if err != nil {
		t.Errorf("Expected nil error for valid book with pointers, got: %v", err)
	}
}

func TestValidate_NilPointer(t *testing.T) {
	var book *ValidBook = nil

	err := Validate(book)
	if err == nil {
		t.Error("Expected error for nil pointer, got nil")
	}

	// Check error type
	if _, ok := err.(*apperrors.ValidateErr); !ok {
		t.Errorf("Expected ValidateErr, got: %T", err)
	}
}

func TestValidate_Nil(t *testing.T) {
	err := Validate(nil)
	if err == nil {
		t.Error("Expected error for nil, got nil")
	}
}

func TestValidate_SQLInjection(t *testing.T) {
	book := ValidBook{
		ID:              1,
		Title:           "SELECT * FROM users; DROP TABLE users;--",
		Genre:           "Programming",
		Author:          "Robert C. Martin",
		PublicationDate: validTime,
	}

	err := Validate(book)
	if err == nil {
		t.Error("Expected error for SQL injection pattern, got nil")
	}
}

func TestValidate_XSSInjection(t *testing.T) {
	book := ValidBook{
		ID:              1,
		Title:           "Clean Code",
		Genre:           "<script>alert('xss')</script>",
		Author:          "Robert C. Martin",
		PublicationDate: validTime,
	}

	err := Validate(book)
	if err == nil {
		t.Error("Expected error for XSS pattern, got nil")
	}
}

func TestValidate_LongString(t *testing.T) {
	longString := "This is a very long string that exceeds the maximum allowed length of 100 characters and should trigger a validation error"

	book := ValidBook{
		ID:              1,
		Title:           longString,
		Genre:           "Programming",
		Author:          "Robert C. Martin",
		PublicationDate: validTime,
	}

	err := Validate(book)
	if err == nil {
		t.Error("Expected error for long string, got nil")
	}
}

func TestValidate_NestedStruct(t *testing.T) {
	nestedBook := NestedBook{
		Book: ValidBook{
			ID:              1,
			Title:           "Clean Code",
			Genre:           "Programming",
			Author:          "Robert C. Martin",
			PublicationDate: validTime,
		},
		CreatedAt: validTime,
	}

	err := Validate(nestedBook)
	if err != nil {
		t.Errorf("Expected nil error for nested struct, got: %v", err)
	}
}

func TestValidate_UnexportedFields(t *testing.T) {
	book := BookWithUnexported{
		ID:     1,
		Title:  "Clean Code",
		Genre:  "Programming",
		Author: "Robert C. Martin",
		// publicationDate is unexported and should be ignored
	}

	err := Validate(book)
	if err != nil {
		t.Errorf("Expected nil error for struct with unexported fields, got: %v", err)
	}
}

func TestValidate_NonStructType(t *testing.T) {
	// Test with non-struct types
	tests := []interface{}{
		"string",
		42,
		[]string{"a", "b"},
		make(map[string]string),
		func() {},
	}

	for i, test := range tests {
		err := Validate(test)
		if err == nil {
			t.Errorf("Test %d: Expected error for non-struct type %T, got nil", i, test)
		}
	}
}

func TestValidate_InvalidStructType(t *testing.T) {
	invalid := InvalidStruct{
		ID:    "not-a-number",
		Title: 42,
	}

	err := Validate(invalid)
	// This should not panic and should handle the type mismatches gracefully
	if err == nil {
		t.Error("Expected error for struct with wrong field types, got nil")
	}
}

func TestValidate_PointerToNilPointer(t *testing.T) {
	var nilBook *ValidBook = nil
	pointerToNil := &nilBook

	err := Validate(pointerToNil)
	if err == nil {
		t.Error("Expected error for pointer to nil pointer, got nil")
	}
}

func TestValidate_EmptyStruct(t *testing.T) {
	empty := struct{}{}

	// Debug: check what reflect sees
	v := reflect.ValueOf(empty)
	t.Logf("Kind: %v, NumField: %d", v.Kind(), v.NumField())

	err := Validate(empty)
	if err == nil {
		t.Error("Expected error for empty struct, got nil")
	} else {
		t.Logf("Got error: %v", err)
	}
}

// Benchmark test to ensure performance doesn't regress
func TestValidate_Performance(t *testing.T) {
	book := ValidBook{
		ID:              1,
		Title:           "Clean Code",
		Genre:           "Programming",
		Author:          "Robert C. Martin",
		PublicationDate: validTime,
	}

	// Run multiple times to ensure no performance regression
	for i := 0; i < 1000; i++ {
		err := Validate(book)
		if err != nil {
			t.Errorf("Unexpected error in performance test: %v", err)
		}
	}
}
