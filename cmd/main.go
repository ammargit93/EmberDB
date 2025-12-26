package main

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var (
	mutx      sync.RWMutex
	DataStore = make(map[string]any)
)

func main() {
	app := fiber.New()

	app.Get("/:name", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello " + c.Params("name"),
		})
	})

	app.Post("/set", SetKey)

	app.Get("/get/:key", GetKey)

	app.Patch("/update", UpdateKey)

	app.Delete("/delete/:key", DeleteKey)

	app.Get("/getall", GetAll)

	app.Listen(":9182")
}
