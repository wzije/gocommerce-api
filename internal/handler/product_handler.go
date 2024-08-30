package handler

import (
	"github.com/ecommerce-api/internal/service"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/exception"
	. "github.com/ecommerce-api/pkg/handler"
	"github.com/ecommerce-api/pkg/http"
	. "github.com/ecommerce-api/pkg/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type productHandler struct {
	service ProductServiceInterface
}

func (b productHandler) List(ctx *fiber.Ctx) error {
	results, err := b.service.List(ctx.Context())

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, results)
}

func (b productHandler) ListWithStock(ctx *fiber.Ctx) error {
	results, err := b.service.ListWithStock(ctx.Context())

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, results)
}

func (b productHandler) GetById(ctx *fiber.Ctx) error {

	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)

	if err != nil {
		return http.JsonError(ctx, exception.ErrInvalidParam)
	}

	result, err := b.service.GetByIDWithStock(ctx.Context(), id)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, result)
}

func NewProductHandler(ProductService ProductServiceInterface) ProductHandlerInterface {
	return &productHandler{service: ProductService}
}

func RegisterProductHandler(db config.DB) ProductHandlerInterface {
	productService := service.NewProductService(db.SqlDB())
	return NewProductHandler(productService)
}
