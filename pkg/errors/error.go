package errors

import (
	"encoding/json"
	"fmt"

	"github.com/govinda-attal/eshop/pkg/errors/codes"
)

// Internal represents internal server error.
func Internal(ee ...error) *Error {
	err := New(codes.ErrInternal, "Internal Server Error")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

// Unauthorized represents an unauthorized request error.
func Unauthorized(ee ...error) *Error {
	err := New(codes.ErrUnauthorized, "Unauthorized")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

// NotFound represents an error when an artifact was not found.
func NotFound(ee ...error) *Error {
	err := New(codes.ErrNotFound, "Not Found")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

// BadRequest represents an invalid request error.
func BadRequest(ee ...error) *Error {
	err := New(codes.ErrBadRequest, "Bad Request")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

// NotImplemented represents a method not implemented error.
func NotImplemented(ee ...error) *Error {
	err := New(codes.ErrNotImplemented, "Not Implemented")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

// ContentTypeNotSupported represents unsupported media type.
func ContentTypeNotSupported(ee ...error) *Error {
	err := New(codes.ErrContentTypeNotSupported, "Unsupported Media Type")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

// StatusConflict represents conflict because of inconsistent or duplicated info.
func StatusConflict(ee ...error) *Error {
	err := New(codes.ErrStatusConflict, "Conflict because of inconsistent or duplicated info")
	if len(ee) > 0 {
		return err.WithError(ee[0])
	}
	return err
}

func ErrCause(err error) error {
	if e, ok := err.(errCauser); ok {
		return e.Cause()
	}
	return nil
}

// Error captures basic information about a status construct.
type Error struct {
	Code    codes.Code `json:"code,omitempty"`
	Message string     `json:"msg,omitempty"`
	Details []*Detail  `json:"details,omitempty"`
	err     error
}

// Detail captures basic information about a status construct.
type Detail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"msg,omitempty"`
}

// WithMessage returns an error status with given message.
func (e *Error) WithMessage(msg string, args ...interface{}) *Error {
	e.Message = fmt.Sprintf(msg, args...)
	return e
}

// WithError returns an error status with given err.Error().
func (e *Error) WithError(err error) *Error {
	e.err = err
	e.Message = fmt.Sprintf("%s: %v", e.Message, err)
	return e
}

// AddDtlMsg returns an error with given message appended to itself.
func (e *Error) WithAdditionalMessages(msgs ...string) *Error {
	for _, m := range msgs {
		d := &Detail{Message: m}
		e.Details = append(e.Details, d)
	}
	return e
}

// WithDetail returns an error with given detail appended to itself.
func (e *Error) WithDetail(code, msg string) *Error {
	d := &Detail{Code: code, Message: msg}
	e.Details = append(e.Details, d)
	return e
}

// New returns a new error with given code and message.
func New(code codes.Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

func (e *Error) Error() string {
	if b, err := json.Marshal(e); err == nil {
		return string(b)
	}
	return `{"code":500, "msg": "error marshal failed"}`
}

// Cause returns an error causer.
func (e *Error) Cause() error {
	return e.err
}

type errCauser interface {
	Cause() error
}

func (e *Error) Is(code codes.Code) bool {
	return e.Code == code
}
