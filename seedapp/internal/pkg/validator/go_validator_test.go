package validator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStruct is a sample struct for validation testing
type TestStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"required,gt=0"`
}

func TestNewGoValidator(t *testing.T) {
	validator := NewGoValidator()
	assert.NotNil(t, validator, "Validator should not be nil")

	// Check if it implements the Validator interface
	var _ Validator = validator // Verify interface implementation at compile time
	assert.NotNil(t, validator, "Validator should not be nil")
}

func TestGoValidator_Validate_Valid(t *testing.T) {
	validator := NewGoValidator()
	ctx := context.Background()

	// Valid test struct
	testData := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	err := validator.Validate(ctx, testData)
	assert.NoError(t, err, "Validation should pass for valid data")
}

func TestGoValidator_Validate_Invalid(t *testing.T) {
	validator := NewGoValidator()
	ctx := context.Background()

	// Invalid test struct - missing required fields
	testData := TestStruct{
		Name:  "",
		Email: "not-an-email",
		Age:   0,
	}

	err := validator.Validate(ctx, testData)
	assert.Error(t, err, "Validation should fail for invalid data")

	// Check if it's an ErrorMap
	errMap, ok := err.(*ErrorMap)
	assert.True(t, ok, "Error should be of type ErrorMap")

	// Check specific field errors
	assert.Contains(t, errMap.Errors, "Name", "Should have error for Name field")
	assert.Contains(t, errMap.Errors, "Email", "Should have error for Email field")
	assert.Contains(t, errMap.Errors, "Age", "Should have error for Age field")

	// Check error messages
	nameErr, exists := errMap.Errors["Name"]
	assert.True(t, exists, "Should have error for Name field")
	assert.Contains(t, nameErr.Error(), "required", "Name error should mention 'required'")

	emailErr, exists := errMap.Errors["Email"]
	assert.True(t, exists, "Should have error for Email field")
	assert.Contains(t, emailErr.Error(), "valid email", "Email error should mention 'valid email'")
}

func TestGoValidator_Validate_InvalidStruct(t *testing.T) {
	validator := NewGoValidator()
	ctx := context.Background()

	// Pass a non-struct value
	err := validator.Validate(ctx, "not a struct")
	assert.Error(t, err, "Should error when validating non-struct")

	// Check if it's NOT an ErrorMap (should be InvalidValidationError)
	_, ok := err.(*ErrorMap)
	assert.False(t, ok, "Error should not be of type ErrorMap")
}
