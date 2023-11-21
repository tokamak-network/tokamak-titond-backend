package api

import "errors"

var (
	ErrBadRequest       = errors.New("bad request")
	ErrResourceNotFound = errors.New("resource not found")
	ErrInternalServer   = errors.New("internal server error")
)
