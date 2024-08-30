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
	"strconv"
)

type inventoryHandler struct {
	service WarehouseServiceInterface
}

func (i inventoryHandler) UpdateWarehouseStatus(ctx *fiber.Ctx) error {
	warehouseID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	isActive := ctx.Query("is_active") == "true"

	if err != nil {
		return http.JsonError(ctx, exception.ErrInvalidParam)
	}

	err = i.service.UpdateWarehouseStatus(ctx.Context(), warehouseID, isActive)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Message: "update warehouse status successfully",
		Data:    nil,
	})
}

func (i inventoryHandler) MyWarehouseList(ctx *fiber.Ctx) error {
	results, err := i.service.MyWarehouseList(ctx.Context())

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.JsonOk(ctx, results)
}

func (i inventoryHandler) IncreaseStock(ctx *fiber.Ctx) error {
	request := new(dto.ChangeStockRequest)

	if err := ctx.BodyParser(request); err != nil {
		return http.JsonError(ctx, err)
	}

	if err := http.ValidateStruct(*request); err != nil {
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Invalid params", err))
	}

	err := i.service.IncreaseStock(ctx.Context(), request.ProductID, request.WarehouseID, request.Quantity)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	//if there is error
	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Message: "increase stock successfully",
		Data:    nil,
	})
}

func (i inventoryHandler) ReduceStock(ctx *fiber.Ctx) error {
	request := new(dto.ChangeStockRequest)

	if err := ctx.BodyParser(request); err != nil {
		return http.JsonError(ctx, err)
	}

	if err := http.ValidateStruct(*request); err != nil {
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Invalid params", err))
	}

	err := i.service.ReduceStock(ctx.Context(), request.ProductID, request.WarehouseID, request.Quantity)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	//if there is error
	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Message: "reduce stock successfully",
		Data:    nil,
	})
}

func (i inventoryHandler) TransferStock(ctx *fiber.Ctx) error {
	request := new(dto.TransferStockRequest)

	if err := ctx.BodyParser(request); err != nil {
		return http.JsonError(ctx, err)
	}

	if err := http.ValidateStruct(*request); err != nil {
		return http.JsonError(ctx, exception.New(fiber.StatusBadRequest, "Invalid params", err))
	}

	err := i.service.TransferStock(
		ctx.Context(),
		request.SourceWarehouseID,
		request.TargetWarehouseID,
		request.ProductID,
		request.Quantity)

	if err != nil {
		return http.JsonError(ctx, err)
	}

	return http.Json(ctx, &http.Response{
		Code:    fiber.StatusCreated,
		Message: "transfer stock successfully",
		Data:    nil,
	})
}

func NewWarehouseHandler(warehouseService WarehouseServiceInterface) WarehouseHandlerInterface {
	return &inventoryHandler{service: warehouseService}
}

func RegisterWarehouseHandler(db config.DB) WarehouseHandlerInterface {
	inventoryService := service.NewWarehouseService(db.SqlDB())
	return NewWarehouseHandler(inventoryService)
}
