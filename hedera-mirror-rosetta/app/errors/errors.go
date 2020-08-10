package errors

import "net/http"

var Errors = map[string]Error{
	BlockNotFound:        New(BlockNotFound, http.StatusBadRequest, true),
	StartMustBeBeforeEnd: New(StartMustBeBeforeEnd, http.StatusBadRequest, false),
}

const (
	BlockNotFound        string = "Block not found"
	StartMustBeBeforeEnd string = "Start must be before end"
)

type Error interface {
	Error() string
	Message() string
	StatusCode() int32
	Retriable() bool
}

type errorStruct struct {
	message    string
	statusCode int32
	retriable  bool
}

func (e *errorStruct) Error() string {
	return e.message
}

func (e *errorStruct) Message() string {
	return e.message
}

func (e *errorStruct) StatusCode() int32 {
	return e.statusCode
}

func (e *errorStruct) Retriable() bool {
	return e.retriable
}

func New(message string, statusCode int32, retriable bool) Error {
	return &errorStruct{
		message:    message,
		statusCode: statusCode,
		retriable:  retriable,
	}
}
