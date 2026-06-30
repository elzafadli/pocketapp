package custommiddleware

import (
	"errors"
	"testing"

	"seedapp/internal/pkg/formatter"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Define custom errors for testing
var (
	ErrNotFound   = errors.New("not found error")
	ErrBadRequest = errors.New("bad request error")
	ErrUnknown    = errors.New("unknown error")
)

func TestGetcode(t *testing.T) {
	// Setup test cases
	codeMap := map[error]formatter.Status{
		ErrNotFound:   formatter.DataNotFound,
		ErrBadRequest: formatter.InvalidRequest,
	}

	tests := []struct {
		name     string
		err      error
		expected formatter.Status
	}{
		{
			name:     "Should return DataNotFound status for not found error",
			err:      ErrNotFound,
			expected: formatter.DataNotFound,
		},
		{
			name:     "Should return InvalidRequest status for bad request error",
			err:      ErrBadRequest,
			expected: formatter.InvalidRequest,
		},
		{
			name:     "Should return InternalServerError status for unknown error",
			err:      ErrUnknown,
			expected: formatter.InternalServerError,
		},
		{
			name:     "Should return InternalServerError status for nil error",
			err:      nil,
			expected: formatter.InternalServerError,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := getcode(tc.err, codeMap)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGethttpstatus(t *testing.T) {
	// Setup test cases
	statusMap := map[error]int{
		ErrNotFound:   fiber.StatusNotFound,
		ErrBadRequest: fiber.StatusBadRequest,
	}

	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "Should return StatusNotFound for not found error",
			err:      ErrNotFound,
			expected: fiber.StatusNotFound,
		},
		{
			name:     "Should return StatusBadRequest for bad request error",
			err:      ErrBadRequest,
			expected: fiber.StatusBadRequest,
		},
		{
			name:     "Should return StatusInternalServerError for unknown error",
			err:      ErrUnknown,
			expected: fiber.StatusInternalServerError,
		},
		{
			name:     "Should return StatusInternalServerError for nil error",
			err:      nil,
			expected: fiber.StatusInternalServerError,
		},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := gethttpstatus(tc.err, statusMap)
			assert.Equal(t, tc.expected, result)
		})
	}
}
