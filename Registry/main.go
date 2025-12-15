package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func WriteToJSON(peer string) {
	file, err := os.OpenFile("peers.txt", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	file.Write([]byte(peer))
	defer file.Close()
}

func main() {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Registry Active",
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		type Response struct {
			PeerIP string `json:"peerip"`
		}
		var resp Response
		resp.PeerIP = c.Context().RemoteAddr().String()
		WriteToJSON(resp.PeerIP)
		return c.JSON(fiber.Map{"peerip": resp.PeerIP})
	})

	app.Listen(":5050")
}
