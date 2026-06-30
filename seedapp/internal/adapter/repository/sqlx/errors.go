package sqlx

import (
	"errors"
	"fmt"
)

var (
	ErrDatabase = errors.New("database error")
	ErrQuery    = errors.New("query error")
	ErrFunction = errors.New("function error")
)

func NewErrDatabase(err error) error {
	return fmt.Errorf("%w: %q", ErrDatabase, err)
}

func NewErrQuery(err error) error {
	return fmt.Errorf("%w: %q", ErrQuery, err)
}
func NewErrFunction(err error) error {
	return fmt.Errorf("%w: %q", ErrFunction, err)
}
