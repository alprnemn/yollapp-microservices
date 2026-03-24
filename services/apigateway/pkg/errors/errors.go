package errors

import "errors"

var (
	ErrNotFound          = errors.New("record not found")
	ErrNotFoundOnRedis   = errors.New("record not found on redis")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicatePhone    = errors.New("phone already exists")
	ErrGeneratePassword  = errors.New("error occurred while generating password")
	ErrPasswordInvalid   = errors.New("invalid password")
	ErrMissingAuthHeader = errors.New("missing Authorization header")
	ErrWrongAuthHeader   = errors.New("wrong Authorization header value")
	ErrSetUserCache      = errors.New("error occurred while setting user to redis")
)
