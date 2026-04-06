package err

import "errors"

var (
	ErrNotFound          = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicatePhone    = errors.New("phone already exists")
	ErrCreateUser        = errors.New("error occurred while creating user")
	ErrGeneratePassword  = errors.New("error generating password")
)
