package identity

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	id := NewID()
	assert.NotEqual(t, uuid.Nil, id.UUID)
	assert.False(t, id.UUID.String() == "00000000-0000-0000-0000-000000000000")
}

func TestNewZeroID(t *testing.T) {
	id := NewZeroID()
	assert.Equal(t, uuid.Nil, id.UUID)
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", id.UUID.String())
}

func TestFromStringOrNil(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid UUID",
			input:    "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expected: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
		},
		{
			name:     "Invalid UUID",
			input:    "invalid-uuid",
			expected: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := FromStringOrNil(tt.input)
			assert.Equal(t, tt.expected, id.UUID.String())
		})
	}
}

func TestListID(t *testing.T) {
	id1 := NewID()
	id2 := NewID()

	list := ListID{id1, id2}

	assert.Equal(t, 2, len(list))
	assert.Equal(t, id1, list[0])
	assert.Equal(t, id2, list[1])
}

func TestNewNullID(t *testing.T) {
	nullID := NewNullID()

	assert.True(t, nullID.Valid)
	assert.NotEqual(t, uuid.Nil, nullID.UUID)
}

func TestNullIDFromStringOrNil(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		valid    bool
	}{
		{
			name:     "Valid UUID",
			input:    "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expected: "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			valid:    true,
		},
		{
			name:     "Invalid UUID",
			input:    "invalid-uuid",
			expected: "00000000-0000-0000-0000-000000000000",
			valid:    true, // The function always sets Valid to true
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "00000000-0000-0000-0000-000000000000",
			valid:    true, // The function always sets Valid to true
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nullID := NullIDFromStringOrNil(tt.input)
			assert.Equal(t, tt.expected, nullID.UUID.String())
			assert.Equal(t, tt.valid, nullID.Valid)
		})
	}
}
