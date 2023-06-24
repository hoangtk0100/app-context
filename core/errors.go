package core

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var ErrInternalServerError = DefaultError{
	StatusField:   http.StatusText(http.StatusInternalServerError),
	ErrorField:    "An internal server error occurred. Please try again later",
	CodeField:     http.StatusInternalServerError,
	GRPCCodeField: codes.Internal,
}

var ErrUnauthorized = DefaultError{
	StatusField:   http.StatusText(http.StatusUnauthorized),
	ErrorField:    "Access denied. Please authenticate with valid credentials",
	CodeField:     http.StatusUnauthorized,
	GRPCCodeField: codes.Unauthenticated,
}

var ErrBadRequest = DefaultError{
	StatusField:   http.StatusText(http.StatusBadRequest),
	ErrorField:    "The request was invalid or contained malformed parameters",
	CodeField:     http.StatusBadRequest,
	GRPCCodeField: codes.InvalidArgument,
}

var ErrNotFound = DefaultError{
	StatusField:   http.StatusText(http.StatusNotFound),
	ErrorField:    "The requested page or resource could not be found",
	CodeField:     http.StatusNotFound,
	GRPCCodeField: codes.NotFound,
}

var ErrForbidden = DefaultError{
	StatusField:   http.StatusText(http.StatusForbidden),
	ErrorField:    "Access to the requested page or resource is forbidden",
	CodeField:     http.StatusForbidden,
	GRPCCodeField: codes.PermissionDenied,
}

var ErrUnsupportedMediaType = DefaultError{
	StatusField:   http.StatusText(http.StatusUnsupportedMediaType),
	ErrorField:    "The media type of the requested resource is not supported",
	CodeField:     http.StatusUnsupportedMediaType,
	GRPCCodeField: codes.InvalidArgument,
}

var ErrConflict = DefaultError{
	StatusField:   http.StatusText(http.StatusConflict),
	ErrorField:    "A conflict occurred with the current state of the resource",
	CodeField:     http.StatusConflict,
	GRPCCodeField: codes.FailedPrecondition,
}

var ErrTimeout = DefaultError{
	StatusField:   http.StatusText(http.StatusRequestTimeout),
	ErrorField:    "The request timed out",
	CodeField:     http.StatusRequestTimeout,
	GRPCCodeField: codes.DeadlineExceeded,
}
