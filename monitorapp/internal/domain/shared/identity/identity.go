// Package identity represents a shared kernel for domain model
package identity

import uuid "github.com/google/uuid"

// ID represents identity of Entity
type ID struct {
	uuid.UUID
}

// String returns the string representation of the ID
func (id ID) String() string {
	return id.UUID.String()
}

// NewID creates a new UUID v7 ID
func NewID() ID {
	_uuid, err := uuid.NewV7()
	if err != nil {
		return NewZeroID()
	}
	return ID{UUID: _uuid}
}

// FromStringOrNil converts string to ID, returning zero ID on parse error
func FromStringOrNil(value string) ID {
	_uuid, err := uuid.Parse(value)
	if err != nil {
		return NewZeroID()
	}
	return ID{UUID: _uuid}
}

// NewZeroID creates a zero (nil) ID
func NewZeroID() ID {
	return ID{UUID: uuid.Nil}
}

// NewDefaultEmptyUUID creates a zero (nil) ID
// This is an alias for NewZeroID for backward compatibility
func NewDefaultEmptyUUID() ID {
	return NewZeroID()
}

// ListID represents a list of IDs
type ListID []ID

// NullID represents a nullable UUID
type NullID struct {
	uuid.NullUUID
}

// NewNullID creates a new NullID with a UUID v7
func NewNullID() NullID {
	_uuid, err := uuid.NewV7()
	if err != nil {
		_uuid = uuid.Nil
	}
	return NullID{
		uuid.NullUUID{
			UUID:  _uuid,
			Valid: true,
		},
	}
}

// NullIDFromStringOrNil converts string to NullID, returning zero UUID with Valid=true on parse error
func NullIDFromStringOrNil(value string) NullID {
	_uuid, err := uuid.Parse(value)
	if err != nil {
		_uuid = uuid.Nil
	}
	return NullID{uuid.NullUUID{
		UUID:  _uuid,
		Valid: true,
	}}
}
