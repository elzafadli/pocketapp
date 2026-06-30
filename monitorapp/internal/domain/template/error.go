package template

import "errors"

var (
	// ErrDataNotFound .
	ErrDataNotFound = errors.New("template: data not found")
	// ErrDataAlreadyExists .
	ErrDataAlreadyExists = errors.New("template: data already exists")
	// ErrDataCreate .
	ErrDataCreate = errors.New("template: data create failed")
	// ErrDataUpdate .
	ErrDataUpdate = errors.New("template: data update failed")
	// ErrDataGet .
	ErrDataGet = errors.New("template: data get failed")
)
