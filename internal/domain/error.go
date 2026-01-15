package domain

import "net/http"

// AppErrorInterface allows the middleware to handle all app errors uniformly
type AppErrorInterface interface {
	error
	GetCode() int
	GetMessage() string
}

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string      { return e.Message }
func (e *AppError) GetCode() int       { return e.Code }
func (e *AppError) GetMessage() string { return e.Message }

func BadRequestError(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

func NotFoundError(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

func UnauthorizedError(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

func ForbiddenError(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

func InternalServerError(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}
