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

type NotFoundError struct {
	Value      string
	Underlying string
}

func NewNotFoundError(e NotFoundError) Error {
	return Error(fmt.Sprintf("%s not found for %s", e.Value, e.Underlying))
}

const (
	ErrNoDatabaseConnection = Error("no database connection")
	ErrSchemaNotCached      = Error("schema not cached")
)
