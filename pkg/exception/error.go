package exception

import (
	"errors"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"net/http"
)

var (
	ErrAuth              = &ApiError{Code: http.StatusUnauthorized, Message: "Invalid Token"}
	ErrNotFound          = &ApiError{Code: http.StatusNotFound, Message: "Not Found"}
	ErrBadRequest        = &ApiError{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrUnprocessedEntity = &ApiError{Code: http.StatusUnprocessableEntity, Message: "Unprocessed Entity"}
	ErrInternalServer    = &ApiError{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
	ErrConflicted        = &ApiError{Code: http.StatusConflict, Message: "Conflicted"}
	ErrParsing           = &ApiError{Code: http.StatusBadRequest, Message: "Parsing data failed"}
	ErrInvalidParam      = &ApiError{Code: http.StatusBadRequest, Message: "Invalid Params"}
	ErrEmptyPassword     = &ApiError{Code: http.StatusBadRequest, Message: "Password Empty"}
	ErrVerification      = &ApiError{Code: http.StatusPreconditionFailed, Message: "Verification Failed"}
	ErrDisableUser       = &ApiError{Code: http.StatusPreconditionFailed, Message: "your account is disabled"}
)

//type ApiError interface {
//	ApiError() (int, string, interface{})
//}

type ApiError struct {
	Code    int
	Message string
	Cause   interface{}
}

func (e ApiError) Error() string {
	return e.Message
}

func (e ApiError) ApiError() (int, string, interface{}) {
	return e.Code, e.Message, e.Cause
}

func New(code int, message string, cause interface{}) *ApiError {
	return &ApiError{Code: code, Message: message, Cause: cause}
}

func Message(message string) error {
	return errors.New(message)
}

type AppErrException struct {
	error
	appError *ApiError
}

func CatchException(err error, appError *ApiError) error {
	return AppErrException{error: err, appError: appError}
}

func TranslateErr(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return CatchException(err, ErrNotFound)
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return CatchException(err, New(http.StatusUnprocessableEntity, "Invalid Transaction", nil))
	case errors.Is(err, gorm.ErrInvalidData):
		return CatchException(err, New(http.StatusNotFound, "Invalid Data", nil))
	case errors.Is(err, err.(*pgconn.PgError)):
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		return CatchException(err, New(http.StatusBadRequest, pgErr.Message, pgErr.Error()))
	default:
		return err
	}
}
