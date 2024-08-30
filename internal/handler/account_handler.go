package handler

import (
	"github.com/ecommerce-api/internal/service"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/dto"
	"github.com/ecommerce-api/pkg/exception"
	. "github.com/ecommerce-api/pkg/handler"
	"github.com/ecommerce-api/pkg/http"
	. "github.com/ecommerce-api/pkg/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type accountHandler struct {
	service AccountServiceInterface
}

func (a accountHandler) Register(ctx *fiber.Ctx) error {
	//init user model
	request := new(dto.AuthRegisterRequest)

	if err := ctx.BodyParser(request); err != nil {
		return http.JsonError(ctx, err)
	}

	if err := http.ValidateStruct(*request); err != nil {
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Invalid params", err))
	}

	result, err := a.service.Register(request)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	//if there is error
	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Data:    result,
		Message: "Success",
	})

}

func (a accountHandler) Login(ctx *fiber.Ctx) error {

	//prepare user
	request := new(dto.AuthAccessTokenRequest)

	//parsing body param to request
	if err := ctx.BodyParser(&request); err != nil {
		return http.JsonError(ctx, exception.ErrInvalidParam)
	}

	//get access token and request
	results, err := a.service.Login(request) //return token and user

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.Json(ctx,
		&http.Response{
			Code:    fiber.StatusOK,
			Message: "login successfully",
			Data:    results,
		})
}

func (a accountHandler) Profile(ctx *fiber.Ctx) error {
	user, err := a.service.Profile()

	//check existing user
	if err != nil {
		logrus.Fatal(err)
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Failed", err))
	}

	//show the profile
	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusOK,
		Message: "Success",
		Data:    user})
}

func (a accountHandler) Logout(ctx *fiber.Ctx) error {
	//TODO implement me
	return nil
}

func NewAccountHandler(accountService AccountServiceInterface) AccountHandlerInterface {
	return &accountHandler{
		service: accountService,
	}
}

func RegisterAccountHandler(db config.DB) AccountHandlerInterface {
	accountService := service.NewAccountService(db)
	return NewAccountHandler(accountService)
}
