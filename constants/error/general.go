package error

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute")
	ErrToManyRequest       = errors.New("too many requests")
	ErrUnauthorize         = errors.New("unauthorize")
	ErrInvalidToken        = errors.New("invalid token")
	ErrForbidden           = errors.New("forbidden")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrToManyRequest,
	ErrUnauthorize,
	ErrInvalidToken,
	ErrForbidden,
}
