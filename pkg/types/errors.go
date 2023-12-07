package types

import "errors"

var (
	ErrBadRequest           = errors.New("bad request")
	ErrResourceNotFound     = errors.New("resource not found")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidComponentType = errors.New("invalid component type")
	ErrComponentDependency  = errors.New("invalid component dependency")
)
