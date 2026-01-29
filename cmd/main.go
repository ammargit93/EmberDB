package main

import (
	"emberdb/internal"
	"emberdb/storage"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Spawn() {
	mu.RLock()
	defer mu.RUnlock()

	val, exists := internal.ArgMap["snapshot"]

	if !exists {
		fmt.Println("No --snapshot flag")
		return
	} else {
		number := val[:len(val)-1]
		unit := val[len(val)-1:]
		duration, exists := internal.DurationMap[unit]
		if !exists {
			fmt.Println("Invalid time unit")
			return
		}
		drtn, _ := strconv.Atoi(number)

		go storage.Snap(time.Duration(drtn) * duration)
	}

}

func main() {
	app := fiber.New()
	internal.Parse(os.Args)

	app.Post("/set", SetKey)

	app.Get("/get/:namespace/:key", GetKey)

	app.Patch("/update", UpdateKey)

	app.Delete("/delete/:namespace/:key", DeleteKey)

	app.Get("/getall", GetAll)

	app.Post("/upload/:namespace/:key", internal.UploadFile)

	Spawn()

	app.Listen(":9182")
}
