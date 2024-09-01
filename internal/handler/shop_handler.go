package handler

import (
	service "github.com/ecommerce-api/internal/service"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/exception"
	. "github.com/ecommerce-api/pkg/handler"
	"github.com/ecommerce-api/pkg/http"
	servicePkg "github.com/ecommerce-api/pkg/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type shopHandler struct {
	shopService servicePkg.ShopServiceInterface
}

func (s shopHandler) List(ctx *fiber.Ctx) error {
	result, err := s.shopService.List(ctx.Context())

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, result)
}

func (s shopHandler) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return http.JsonError(ctx, exception.ErrInvalidParam)
	}

	result, err := s.shopService.GetByID(ctx.Context(), id)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, result)

}

func (s shopHandler) GetProducts(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return http.JsonError(ctx, exception.ErrInvalidParam)
	}

	result, err := s.shopService.Products(ctx.Context(), id)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, result)

}

func (s shopHandler) GetOrders(ctx *fiber.Ctx) error {

	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return http.JsonError(ctx, exception.ErrInvalidParam)
	}

	result, err := s.shopService.Orders(ctx.Context(), id)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, result)
}

func NewShopHandler(shopService servicePkg.ShopServiceInterface) ShopHandlerInterface {
	return &shopHandler{shopService: shopService}
}

func RegisterShopHandler(db config.DB) ShopHandlerInterface {
	shopService := service.NewShopService(db.SqlDB())
	return NewShopHandler(shopService)
}
