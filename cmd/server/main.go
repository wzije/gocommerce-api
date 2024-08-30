package main

import (
	"fmt"
	"github.com/ecommerce-api/internal/router"
	"github.com/ecommerce-api/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
)

func init() {
	log.Println("Execute Init...")
	config.Load()
	//tasks.RegisterTask()
}

func listen() {
	app := fiber.New()
	app.Use(logger.New())

	longTask := new(sync.WaitGroup)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello GoEcommerce API")
	})

	//register router
	db := config.SqlDBLoad()
	router.RegisterApiRouter(app, db, longTask)

	go log.Fatal(app.Listen(fmt.Sprintf("0.0.0.0:%s", config.Config.AppPort)))

	longTask.Wait()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
}

func main() {
	fmt.Printf("Listening [::]:%s\n", config.Config.AppPort)
	listen()
}
