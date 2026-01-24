package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/set", SetKey)

	app.Get("/get/:namespace/:key", GetKey)

	// app.Post("/mset", MSet)

	// app.Post("/mget", MGet)

	app.Patch("/update", UpdateKey)

	app.Delete("/delete/:namespace/:key", DeleteKey)

	app.Get("/getall", GetAll)

	// file logic
	// app.Post("/upload/:key", internal.UploadFile)

	app.Listen(":9182")
}
