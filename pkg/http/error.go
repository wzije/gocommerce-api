package http

import (
	"errors"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

var (
	ErrEmailNotFound      = errors.New("email not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrEmptyPassword      = errors.New("password can't be empty")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidAuthToken   = errors.New("invalid auth-token")
	ErrInvalidCredentials = errors.New("invalid password credentials")
	ErrUnauthorized       = errors.New("unauthorized")
)

type ErrorResponse struct {
	Field string `json:"field,omitempty"`
	Tag   string `json:"tag,omitempty"`
	Value string `json:"value,omitempty"`
}

func ErrorHttp(err error) *Response {
	logrus.Error(err.Error())
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return &Response{Code: http.StatusNotFound, Message: err.Error()}
	default:
		return &Response{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}

var validate = validator.New()

func ValidateStruct(data interface{}) []*ErrorResponse {
	var errs []*ErrorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = helper.ToSnakeCase(err.Field())
			element.Tag = helper.ToSnakeCase(err.Tag())
			element.Value = helper.ToSnakeCase(err.Param())
			errs = append(errs, &element)
		}
	}
	return errs
}
