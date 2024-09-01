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
	"sync"
)

type orderHandler struct {
	service OrderServiceInterface
}

func (o orderHandler) MyCustomerOrders(ctx *fiber.Ctx) error {
	orders, err := o.service.MyOrders(ctx.Context())

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusOK,
		Message: "fetch my orders",
		Data:    orders,
	})
}

func (o orderHandler) PaymentOrder(ctx *fiber.Ctx) error {
	request := new(dto.PaymentRequest)

	if err := ctx.BodyParser(request); err != nil {
		return http.JsonError(ctx, err)
	}

	if err := http.ValidateStruct(*request); err != nil {
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Invalid params", err))
	}

	err := o.service.PaymentOrder(ctx.Context(), request)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	//if there is error
	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Message: "payment order created",
		Data:    nil,
	})
}

func (o orderHandler) CheckoutOrder(ctx *fiber.Ctx) error {
	request := new(dto.OrderRequest)

	if err := ctx.BodyParser(request); err != nil {
		return http.JsonError(ctx, err)
	}

	if err := http.ValidateStruct(*request); err != nil {
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Invalid params", err))
	}

	order, err := o.service.CheckoutOrder(ctx.Context(), request)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	//if there is error
	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Message: "order created",
		Data:    order,
	})
}

func NewOrderHandler(orderService OrderServiceInterface) OrderHandlerInterface {
	return &orderHandler{service: orderService}
}

func RegisterOrderHandler(db config.DB, task *sync.WaitGroup) OrderHandlerInterface {
	orderService := service.NewOrderService(db.SqlDB(), task)
	return NewOrderHandler(orderService)
}
