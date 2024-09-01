package handler

import "github.com/gofiber/fiber/v2"

type AccountHandlerInterface interface {
	Profile(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
}

type UserHandlerInterface interface {
	BaseHandlerInterface
}

type ShopHandlerInterface interface {
	List(ctx *fiber.Ctx) error
	GetByID(ctx *fiber.Ctx) error
	GetProducts(ctx *fiber.Ctx) error
	GetOrders(ctx *fiber.Ctx) error
}

type ProductHandlerInterface interface {
	List(ctx *fiber.Ctx) error
	ListWithStock(ctx *fiber.Ctx) error
	GetById(ctx *fiber.Ctx) error
}

type OrderHandlerInterface interface {
	MyCustomerOrders(ctx *fiber.Ctx) error
	CheckoutOrder(ctx *fiber.Ctx) error
	PaymentOrder(ctx *fiber.Ctx) error
}

type WarehouseHandlerInterface interface {
	MyWarehouseList(ctx *fiber.Ctx) error
	CreateProductInventory(ctx *fiber.Ctx) error
	UpdateWarehouseStatus(ctx *fiber.Ctx) error
	IncreaseStock(ctx *fiber.Ctx) error
	ReduceStock(ctx *fiber.Ctx) error
	TransferStock(ctx *fiber.Ctx) error
}
