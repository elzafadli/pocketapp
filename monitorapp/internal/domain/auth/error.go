package auth

import "errors"

var (
	ErrInvalidBasicAuth = errors.New("auth: invalid basic auth")
	ErrInvalidApiKey    = errors.New("auth: invalid api key")
)
