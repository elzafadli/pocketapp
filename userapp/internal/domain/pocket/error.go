package pocket

import "errors"

const (
	ErrCodePocketNotFound = "POCKET_NOT_FOUND"
	ErrMsgPocketNotFound  = "Pocket item not found"
)

var (
	ErrPocketNotFound    = errors.New("pocket item not found")
	ErrDataAlreadyExists = errors.New("pocket: data already exists")
	ErrDataCreate        = errors.New("pocket: failed to create")
	ErrDataUpdate        = errors.New("pocket: failed to update")
	ErrDataDelete        = errors.New("pocket: failed to delete")
	ErrDataGet           = errors.New("pocket: failed to get")
)
