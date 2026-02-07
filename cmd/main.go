package main

import (
	"emberdb/internal"
	"emberdb/storage"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	internal.Parse(os.Args)
	err := storage.LoadFromJSON()
	if err != nil {
		fmt.Println("Invalid snapshot", err)
		storage.ReplayWAL()
	}

	app.Post("/set", SetKey)

	app.Get("/get/:namespace/:key", GetKey)

	app.Patch("/update", UpdateKey)

	app.Delete("/delete/:namespace/:key", DeleteKey)

	app.Get("/getall", GetAll)

	app.Post("/upload/:namespace/:key", internal.UploadFile)

	storage.Spawn()

	app.Listen(":9182")
}
