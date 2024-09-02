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
	route := a.router.Use(
		middleware.AuthAPIMiddleware, func(ctx *fiber.Ctx) error {
			security.ParsePayload(ctx)
			return ctx.Next()
		})

	accountHandler := handler.RegisterAccountHandler(a.db)
	route.Get("me", accountHandler.Profile)

	// shop route
	shopHandler := handler.RegisterShopHandler(a.db)
	route.Get("/shops", shopHandler.List)
	route.Get("/shops/:id", shopHandler.GetByID)
	route.Get("/shops/:id/products", shopHandler.GetProducts)
	route.Get("/shops/:id/orders", shopHandler.GetOrders)

	//product route
	productHandler := handler.RegisterProductHandler(a.db)
	route.Get("/products", productHandler.ListWithStock)
	route.Get("/products/:id", productHandler.GetById)

	//order route
	orderHandler := handler.RegisterOrderHandler(a.db, a.task)
	route.Get("/orders", orderHandler.MyCustomerOrders)
	route.Post("/orders/checkout", orderHandler.CheckoutOrder)
	route.Post("/orders/payment", orderHandler.PaymentOrder)

	//warehouse route
	warehouseHandler := handler.RegisterWarehouseHandler(a.db)
	route.Get("/warehouses", warehouseHandler.MyWarehouseList)
	route.Get("/warehouses/:id", warehouseHandler.MyWarehouseByID)
	route.Post("/warehouses/create", warehouseHandler.CreateWarehouse)
	route.Post("/warehouses/inventories/create", warehouseHandler.CreateProductInventory)
	route.Post("/warehouses/inventories/increase", warehouseHandler.IncreaseStock)
	route.Post("/warehouses/inventories/reduce", warehouseHandler.ReduceStock)
	route.Post("/warehouses/inventories/transfer", warehouseHandler.TransferStock)
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
