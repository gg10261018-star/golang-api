package errorx

import "fmt"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

const (
	CodeBadRequest          = 400
	CodeUnauthorized        = 401
	CodeConflict            = 409
	CodeInternalServerError = 500
)

func New(code int, msg ...any) *AppError {
	return &AppError{
		Code:    code,
		Message: fmt.Sprint(msg...),
	}
}
