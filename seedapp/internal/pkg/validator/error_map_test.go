package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorMap(t *testing.T) {
	// Create a map of errors
	errs := map[string]error{
		"field1": errors.New("error1"),
		"field2": errors.New("error2"),
	}

	// Create a new ErrorMap
	err := NewErrorMap(errs)
	assert.NotNil(t, err, "ErrorMap should not be nil")

	// Check if it's an ErrorMap
	errMap, ok := err.(*ErrorMap)
	assert.True(t, ok, "Error should be of type ErrorMap")
	assert.Equal(t, errs, errMap.Errors, "ErrorMap should contain the provided errors")
}

func TestErrorMap_Error(t *testing.T) {
	// Create a map of errors
	errs := map[string]error{
		"field1": errors.New("error1"),
		"field2": errors.New("error2"),
	}

	// Create a new ErrorMap
	errMap := &ErrorMap{Errors: errs}

	// Get the error message
	errMsg := errMap.Error()

	// The error message should contain all field:error pairs
	assert.Contains(t, errMsg, "field1:error1", "Error message should contain field1:error1")
	assert.Contains(t, errMsg, "field2:error2", "Error message should contain field2:error2")
	assert.Contains(t, errMsg, ";", "Error message should contain semicolon separator")
}

func TestErrorResponse(t *testing.T) {
	// Test the ErrorResponse struct
	resp := ErrorResponse{
		Field: "username",
		Tag:   "required",
		Value: "test",
	}

	if resp.Field != "username" {
		t.Errorf("Expected Field to be 'username', got '%s'", resp.Field)
	}
	if resp.Tag != "required" {
		t.Errorf("Expected Tag to be 'required', got '%s'", resp.Tag)
	}
	if resp.Value != "test" {
		t.Errorf("Expected Value to be 'test', got '%s'", resp.Value)
	}
}
