package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	httpLogger "github.com/gofiber/fiber/v2/middleware/logger"
	fiber_redis "github.com/gofiber/storage/redis"
)

func main() {
	app := fiber.New()
	app.Use(requestid.New())
	loggerFormat := fmt.Sprintf("demo-request-%s ${pid} ${locals:requestid} ${status} - ${method} ${path} \n", "local")
	app.Use(httpLogger.New(httpLogger.Config{
		Format:   loggerFormat,
		TimeZone: "Asia/Kolkata",
	}))

	redisStoage := fiber_redis.New(fiber_redis.Config{
		URL:   "redis://localhost:6379",
		Reset: false,
	})
	app.Use(cache.New(cache.Config{
		Expiration: 15 * time.Second,
		Storage:    redisStoage,
		Methods:    []string{fiber.MethodGet, fiber.MethodHead, fiber.MethodPost},
	}))

	app.Get("/user/:name", func(c *fiber.Ctx) error {
		time.Sleep(2 * time.Second)
		return c.SendString("Hello, " + c.Params("name"))
	})

	app.Get("/foo", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3001")
}
