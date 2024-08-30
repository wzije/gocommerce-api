package handler

import "github.com/gofiber/fiber/v2"

type AccountHandlerInterface interface {
	Profile(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}
