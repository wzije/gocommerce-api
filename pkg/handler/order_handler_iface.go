package handler

import "github.com/gofiber/fiber/v2"

type OrderHandlerInterface interface {
	MyListOrder(ctx *fiber.Ctx) error
	CheckoutOrder(ctx *fiber.Ctx) error
	PaymentOrder(ctx *fiber.Ctx) error
}
