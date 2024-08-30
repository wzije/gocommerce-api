package handler

import "github.com/gofiber/fiber/v2"

type WarehouseHandlerInterface interface {
	MyWarehouseList(ctx *fiber.Ctx) error
	UpdateWarehouseStatus(ctx *fiber.Ctx) error
	IncreaseStock(ctx *fiber.Ctx) error
	ReduceStock(ctx *fiber.Ctx) error
	TransferStock(ctx *fiber.Ctx) error
}
