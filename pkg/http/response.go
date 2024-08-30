package http

import (
	"errors"
	appErrors "github.com/ecommerce-api/pkg/exception"
	"github.com/ecommerce-api/pkg/helper"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type BodyRemote map[string]interface{}

type Response struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func Json(ctx *fiber.Ctx, response *Response) error {
	if response.Code == 0 {
		response.Code = fiber.StatusBadRequest
	}

	return ctx.Status(response.Code).JSON(response)
}

func JsonOk(ctx *fiber.Ctx, data interface{}) error {
	response := Response{Code: fiber.StatusOK, Message: "Success", Data: data}
	return ctx.Status(response.Code).JSON(response)
}

func JsonCreated(ctx *fiber.Ctx, data interface{}) error {
	response := Response{Code: fiber.StatusCreated, Message: "Created", Data: data}
	return ctx.Status(response.Code).JSON(response)
}

func JsonError(ctx *fiber.Ctx, err error) error {
	var apiErr appErrors.ApiError
	if errors.As(err, &apiErr) {
		log.Errorf("Code=%d, message=%s, cause=%s", apiErr.Code, apiErr.Message, apiErr.Cause)
		return ctx.Status(apiErr.Code).
			JSON(Response{Code: apiErr.Code, Message: apiErr.Message, Data: apiErr.Cause})
	} else {
		log.Errorf(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(Response{Code: fiber.StatusInternalServerError, Message: err.Error(), Data: err})
	}
}

func JsonParseBody(ctx *fiber.Ctx, bodyRemote *BodyRemote) error {
	body := *bodyRemote

	code, err := helper.AnyToInt(body["code"])
	if err != nil {
		return JsonError(ctx, err)
	}

	if *code > 500 {
		return ctx.Status(500).JSON(bodyRemote)
	}

	return ctx.Status(*code).JSON(bodyRemote)
}
