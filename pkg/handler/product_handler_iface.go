package handler

import "github.com/gofiber/fiber/v2"

type ProductHandlerInterface interface {
	List(ctx *fiber.Ctx) error
	ListWithStock(ctx *fiber.Ctx) error
	GetById(ctx *fiber.Ctx) error
}
