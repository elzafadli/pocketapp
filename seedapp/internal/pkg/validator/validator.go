package validator

import "context"

//go:generate mockgen -source=validator.go -destination=mocks/validator_mock.go -package=mocks
type Validator interface {
	Validate(ctx context.Context, data interface{}) error
}
