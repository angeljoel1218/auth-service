package apierror

import (
	"errors"
	"net/http"
)

type Type string

const (
	Unauthorized          Type = "UNAUTHORIZED"
	BadRequest            Type = "BAD_REQUEST"
	Conflict              Type = "CONFLICT"
	InternalServer        Type = "INTERNAL"
	NotFound              Type = "NOT_FOUND"
	RequestEntityTooLarge Type = "PAYLOAD_TOO_LARGE"
	ServiceUnavailable    Type = "SERVICE_UNAVAILABLE"
	UnsupportedMediaType  Type = "UNSUPPORTED_MEDIA_TYPE"
	Forbidden             Type = "FORBIDDEN"
	AppError              Type = "ERROR_APPLICATION"
)

type ErrorApi struct {
	Type    Type   `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorApi) Error() string {
	return e.Message
}

func NewErrorApi(status Type, message string) *ErrorApi {
	return &ErrorApi{
		Type:    status,
		Message: message,
	}
}

func NewErrorLogic(status Type, code int, message string) *ErrorApi {
	return &ErrorApi{
		Type:    status,
		Code:    code,
		Message: message,
	}
}

func (e *ErrorApi) Status() int {
	switch e.Type {
	case Unauthorized:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case InternalServer:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case RequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case AppError:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}

func HttpStatus(err error) int {
	var e *ErrorApi
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}
