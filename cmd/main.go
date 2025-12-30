package main

import (
	"emberdb/internal"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/:name", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"message": "Hello " + c.Params("name")}) })

	app.Post("/set", SetKey)

	app.Get("/get/:key", GetKey)

	app.Post("/mset", MSet)

	app.Post("/mget", MGet)

	app.Patch("/update", UpdateKey)

	app.Delete("/delete/:key", DeleteKey)

	app.Get("/getall", GetAll)

	// file logic
	app.Post("/upload/:key", internal.UploadFile)

	app.Listen(":9182")
}
