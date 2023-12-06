package errors

import "fmt"

// Error is the custom error type.
type Error string

func (e Error) Error() string {
	return string(e)
}

type BadRequestError struct {
	Expected   string
	Got        string
	Underlying string
}

func NewBadRequestError(e BadRequestError) Error {
	return Error(fmt.Sprintf("%s: expected %s, got %s", e.Underlying, e.Expected, e.Got))
}

const (
	ErrNoDatabaseConnection = Error("no database connection")
)
