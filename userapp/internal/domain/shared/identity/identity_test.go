package identity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type IdentityTestSuite struct {
	suite.Suite
}

func TestIdentityTestSuite(t *testing.T) {
	suite.Run(t, new(IdentityTestSuite))
}

func (suite *IdentityTestSuite) TestID_String() {
	// Test with a valid UUID
	validUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	id := ID{UUID: validUUID}

	result := id.String()

	suite.Equal("550e8400-e29b-41d4-a716-446655440000", result)
}

func (suite *IdentityTestSuite) TestID_String_ZeroID() {
	// Test with zero UUID
	zeroID := NewZeroID()

	result := zeroID.String()

	suite.Equal("00000000-0000-0000-0000-000000000000", result)
}

func (suite *IdentityTestSuite) TestNewID() {
	// Test that NewID creates a valid UUID v7
	id := NewID()

	// Verify it's not a zero UUID
	suite.NotEqual(uuid.Nil, id.UUID)
	// Verify it's a valid UUID
	suite.NotEmpty(id.String())
	// Verify it's a UUID v7 (version 7)
	suite.Equal(byte(7), id.UUID[6]>>4)
}

func (suite *IdentityTestSuite) TestNewID_Uniqueness() {
	// Test that NewID creates unique IDs
	id1 := NewID()
	id2 := NewID()

	suite.NotEqual(id1.UUID, id2.UUID)
	suite.NotEqual(id1.String(), id2.String())
}

func (suite *IdentityTestSuite) TestFromStringOrNil_ValidUUID() {
	// Test with valid UUID string
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"

	id := FromStringOrNil(validUUIDStr)

	suite.Equal(validUUIDStr, id.String())
	suite.NotEqual(uuid.Nil, id.UUID)
}

func (suite *IdentityTestSuite) TestFromStringOrNil_InvalidUUID() {
	// Test with invalid UUID string
	invalidUUIDStr := "invalid-uuid-string"

	id := FromStringOrNil(invalidUUIDStr)

	// Should return zero ID
	suite.Equal(uuid.Nil, id.UUID)
	suite.Equal("00000000-0000-0000-0000-000000000000", id.String())
}

func (suite *IdentityTestSuite) TestFromStringOrNil_EmptyString() {
	// Test with empty string
	id := FromStringOrNil("")

	// Should return zero ID
	suite.Equal(uuid.Nil, id.UUID)
	suite.Equal("00000000-0000-0000-0000-000000000000", id.String())
}

func (suite *IdentityTestSuite) TestFromStringOrNil_MalformedUUID() {
	// Test with malformed UUID string
	malformedUUIDStr := "not-a-uuid"

	id := FromStringOrNil(malformedUUIDStr)

	// Should return zero ID
	suite.Equal(uuid.Nil, id.UUID)
}

func (suite *IdentityTestSuite) TestFromStringOrNil_UUIDV4() {
	// Test with UUID v4 string
	uuidV4Str := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	id := FromStringOrNil(uuidV4Str)

	suite.Equal(uuidV4Str, id.String())
	suite.NotEqual(uuid.Nil, id.UUID)
}

func (suite *IdentityTestSuite) TestNewZeroID() {
	// Test that NewZeroID returns zero UUID
	id := NewZeroID()

	suite.Equal(uuid.Nil, id.UUID)
	suite.Equal("00000000-0000-0000-0000-000000000000", id.String())
}

func (suite *IdentityTestSuite) TestNewZeroID_MultipleCalls() {
	// Test that multiple calls return the same zero ID
	id1 := NewZeroID()
	id2 := NewZeroID()

	suite.Equal(id1.UUID, id2.UUID)
	suite.Equal(uuid.Nil, id1.UUID)
	suite.Equal(uuid.Nil, id2.UUID)
}

func (suite *IdentityTestSuite) TestNewDefaultEmptyUUID() {
	// Test that NewDefaultEmptyUUID returns zero UUID
	id := NewDefaultEmptyUUID()

	suite.Equal(uuid.Nil, id.UUID)
	suite.Equal("00000000-0000-0000-0000-000000000000", id.String())
}

func (suite *IdentityTestSuite) TestNewDefaultEmptyUUID_IsAlias() {
	// Test that NewDefaultEmptyUUID is an alias for NewZeroID
	id1 := NewDefaultEmptyUUID()
	id2 := NewZeroID()

	suite.Equal(id1.UUID, id2.UUID)
	suite.Equal(uuid.Nil, id1.UUID)
}

func (suite *IdentityTestSuite) TestListID() {
	// Test that ListID can hold multiple IDs
	id1 := NewID()
	id2 := NewID()
	id3 := NewID()

	list := ListID{id1, id2, id3}

	suite.Len(list, 3)
	suite.Equal(id1, list[0])
	suite.Equal(id2, list[1])
	suite.Equal(id3, list[2])
}

func (suite *IdentityTestSuite) TestListID_Empty() {
	// Test empty ListID
	list := ListID{}

	suite.Len(list, 0)
	suite.NotNil(list)
}

func (suite *IdentityTestSuite) TestNewNullID() {
	// Test that NewNullID creates a valid NullID with UUID v7
	nullID := NewNullID()

	suite.True(nullID.Valid)
	suite.NotEqual(uuid.Nil, nullID.UUID)
	// Verify it's a UUID v7 (version 7)
	suite.Equal(byte(7), nullID.UUID[6]>>4)
}

func (suite *IdentityTestSuite) TestNewNullID_Uniqueness() {
	// Test that NewNullID creates unique IDs
	nullID1 := NewNullID()
	nullID2 := NewNullID()

	suite.NotEqual(nullID1.UUID, nullID2.UUID)
	suite.True(nullID1.Valid)
	suite.True(nullID2.Valid)
}

func (suite *IdentityTestSuite) TestNullIDFromStringOrNil_ValidUUID() {
	// Test with valid UUID string
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"

	nullID := NullIDFromStringOrNil(validUUIDStr)

	suite.True(nullID.Valid)
	suite.Equal(validUUIDStr, nullID.UUID.String())
	suite.NotEqual(uuid.Nil, nullID.UUID)
}

func (suite *IdentityTestSuite) TestNullIDFromStringOrNil_InvalidUUID() {
	// Test with invalid UUID string
	invalidUUIDStr := "invalid-uuid-string"

	nullID := NullIDFromStringOrNil(invalidUUIDStr)

	// Should return zero UUID but Valid=true
	suite.True(nullID.Valid)
	suite.Equal(uuid.Nil, nullID.UUID)
}

func (suite *IdentityTestSuite) TestNullIDFromStringOrNil_EmptyString() {
	// Test with empty string
	nullID := NullIDFromStringOrNil("")

	// Should return zero UUID but Valid=true
	suite.True(nullID.Valid)
	suite.Equal(uuid.Nil, nullID.UUID)
}

func (suite *IdentityTestSuite) TestNullIDFromStringOrNil_MalformedUUID() {
	// Test with malformed UUID string
	malformedUUIDStr := "not-a-uuid"

	nullID := NullIDFromStringOrNil(malformedUUIDStr)

	// Should return zero UUID but Valid=true
	suite.True(nullID.Valid)
	suite.Equal(uuid.Nil, nullID.UUID)
}

func (suite *IdentityTestSuite) TestNullIDFromStringOrNil_UUIDV4() {
	// Test with UUID v4 string
	uuidV4Str := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	nullID := NullIDFromStringOrNil(uuidV4Str)

	suite.True(nullID.Valid)
	suite.Equal(uuidV4Str, nullID.UUID.String())
	suite.NotEqual(uuid.Nil, nullID.UUID)
}

// Additional edge case tests

func (suite *IdentityTestSuite) TestID_String_WithDifferentVersions() {
	// Test String() with different UUID versions
	testCases := []struct {
		name     string
		uuidStr  string
		expected string
	}{
		{"UUID v1", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		{"UUID v4", "f47ac10b-58cc-4372-a567-0e02b2c3d479", "f47ac10b-58cc-4372-a567-0e02b2c3d479"},
		{"UUID v7", "01890a5d-a756-774b-a895-353aef2acd3f", "01890a5d-a756-774b-a895-353aef2acd3f"},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			uuid := uuid.MustParse(tc.uuidStr)
			id := ID{UUID: uuid}
			suite.Equal(tc.expected, id.String())
		})
	}
}

func (suite *IdentityTestSuite) TestFromStringOrNil_WithUUIDV7() {
	// Test FromStringOrNil with UUID v7 string
	uuidV7Str := "01890a5d-a756-774b-a895-353aef2acd3f"

	id := FromStringOrNil(uuidV7Str)

	suite.Equal(uuidV7Str, id.String())
	suite.NotEqual(uuid.Nil, id.UUID)
}

func (suite *IdentityTestSuite) TestListID_Append() {
	// Test appending to ListID
	list := ListID{}
	id1 := NewID()
	id2 := NewID()

	list = append(list, id1, id2)

	suite.Len(list, 2)
	suite.Equal(id1, list[0])
	suite.Equal(id2, list[1])
}

func (suite *IdentityTestSuite) TestNullID_String() {
	// Test NullID string representation
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"
	nullID := NullIDFromStringOrNil(validUUIDStr)

	suite.Equal(validUUIDStr, nullID.UUID.String())
}

func (suite *IdentityTestSuite) TestNullID_ZeroValue() {
	// Test zero value NullID
	var nullID NullID

	suite.False(nullID.Valid)
	suite.Equal(uuid.Nil, nullID.UUID)
}

// Benchmark tests (optional but good for completeness)
func BenchmarkNewID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewID()
	}
}

func BenchmarkFromStringOrNil(b *testing.B) {
	uuidStr := "550e8400-e29b-41d4-a716-446655440000"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FromStringOrNil(uuidStr)
	}
}

func BenchmarkID_String(b *testing.B) {
	id := NewID()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = id.String()
	}
}
