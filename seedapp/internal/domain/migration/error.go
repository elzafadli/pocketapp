package migration

import (
	"errors"
	"fmt"
)

var (
	ErrSchemaNotFound    = errors.New("schema not found")
	ErrMigrationNotFound = errors.New("migration not found")
	ErrFailedMigrate     = errors.New("failed to migrate")
)

func NewErrFailedMigrate(errorMessage string) error {
	return fmt.Errorf("%w: %s", ErrFailedMigrate, errorMessage)
}
