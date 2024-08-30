package router

import (
	"github.com/ecommerce-api/internal/handler"
	"github.com/ecommerce-api/pkg/config"
	"github.com/ecommerce-api/pkg/middleware"
	"github.com/ecommerce-api/pkg/router"
	"github.com/ecommerce-api/pkg/security"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"sync"
)

type apiRouter struct {
	router fiber.Router
	db     config.DB
	task   *sync.WaitGroup
}

func (a apiRouter) GuestRouter() {

	a.router.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	accountHandler := handler.RegisterAccountHandler(a.db)
	a.router.Post("login", accountHandler.Login)
	a.router.Post("register", accountHandler.Register)
}

func (a apiRouter) AuthRouter() {
	//setup auth middleware
	route := a.router.Group("", middleware.AuthAPIMiddleware, func(ctx *fiber.Ctx) error {
		security.ParsePayload(ctx)
		return ctx.Next()
	})

	accountHandler := handler.RegisterAccountHandler(a.db)
	route.Get("me", accountHandler.Profile)

	//product route
	productHandler := handler.RegisterProductHandler(a.db)
	route.Get("/products", productHandler.ListWithStock)
	route.Get("/products/:id", productHandler.GetById)

	orderHandler := handler.RegisterOrderHandler(a.db, a.task)
	route.Get("/orders", orderHandler.MyListOrder)
	route.Post("/orders/checkout", orderHandler.CheckoutOrder)
	route.Post("/orders/payment", orderHandler.PaymentOrder)

	warehouseHandler := handler.RegisterWarehouseHandler(a.db)
	route.Get("/warehouses", warehouseHandler.MyWarehouseList)
	route.Post("/warehouses/increase", warehouseHandler.IncreaseStock)
	route.Post("/warehouses/reduce", warehouseHandler.ReduceStock)
	route.Post("/warehouses/transfer", warehouseHandler.TransferStock)
	route.Post("/warehouses/:id/status", warehouseHandler.UpdateWarehouseStatus)

}

func NewApiRouter(r fiber.Router, db config.DB, task *sync.WaitGroup) router.AppRouterInterface {
	return &apiRouter{router: r, db: db, task: task}
}

func RegisterApiRouter(r fiber.Router, db config.DB, task *sync.WaitGroup) {
	api := r.Group("/api", func(c *fiber.Ctx) error {
		return c.Next()
	})
	route := NewApiRouter(api, db, task)
	r.Use(cors.New(cors.ConfigDefault))
	route.GuestRouter()
	route.AuthRouter()
}
