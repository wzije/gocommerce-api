package handler

import (
	"github.com/gofiber/fiber/v2"
)

type BaseHandlerInterface interface {
	Fetch(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	ById(ctx *fiber.Ctx) error
	Store(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}
