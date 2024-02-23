package main

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	datastore "github.com/yohansp/go-biometric-api/internal/datastore/sql"
	"github.com/yohansp/go-biometric-api/internal/handlers"
)

func test(c *fiber.Ctx) error {
	fmt.Println("middleware")
	return c.Next()
}

func main() {

	// database
	fmt.Println("starting...")
	datastore.InitDb()

	// setup web container
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			fmt.Println("Error handler start, ", err)
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			return c.Status(code).JSON(err)
		},
	})

	// to fiber handle the panic
	app.Use(recover.New())

	// authorization (permit api)
	handlers.HandlerAuthorizationRoute(app)

	// settings (no need permit api)
	handlers.HandlerSettingRoute(app)

	// admin route
	handlers.HandlerAdminRoute(app)

	app.Listen(":8181")
}
