package main

import (
	"emberdb/internal"
	"emberdb/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/set", SetKey)

	app.Get("/get/:namespace/:key", GetKey)

	app.Patch("/update", UpdateKey)

	app.Delete("/delete/:namespace/:key", DeleteKey)

	app.Get("/getall", GetAll)

	app.Post("/upload/:namespace/:key", internal.UploadFile)

	go storage.Snap("data/snapshot.json")

	app.Listen(":9182")
}
