package error

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPasswordIncorrect = errors.New("password incorrect")
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrPhoneExists       = errors.New("phone already exists")
	ErrPasswordDoesMatch = errors.New("password does not match")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordIncorrect,
	ErrUsernameExists,
	ErrEmailExists,
	ErrPhoneExists,
	ErrPasswordDoesMatch,
}
