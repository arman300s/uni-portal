package contracts

import "errors"

var (
	ErrEmailInUse         = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrRoleNotFound       = errors.New("role not found")
	ErrSubjectNotFound    = errors.New("subject not found")
	ErrForbidden          = errors.New("forbidden")
)
