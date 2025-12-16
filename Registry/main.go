package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func WriteToFile(peer string) {
	file, err := os.OpenFile("peers.txt", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	file.Write([]byte(peer + "\n"))
	defer file.Close()
}

func peerExists(addr string) bool {
	content, err := os.ReadFile("peers.txt")
	if err != nil {
		log.Fatalln(err)
	}
	peers := strings.Split(string(content), "\n")
	peers = peers[:len(peers)-1]
	return slices.Contains(peers, addr)

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
		var ip string
		localAddr := strings.Split(c.Context().LocalAddr().String(), ":")
		if len(localAddr) > 1 {
			ip = localAddr[0]
		}
		resp.PeerIP = ip + ":" + c.Get("X-Port")
		fmt.Println(resp.PeerIP)
		if !peerExists(resp.PeerIP) {
			WriteToFile(resp.PeerIP)
		}
		return c.JSON(fiber.Map{"peerip": resp.PeerIP})
	})

	app.Get("/fetch-peers", func(c *fiber.Ctx) error {
		content, err := os.ReadFile("peers.txt")
		if err != nil {
			log.Fatalln(err)
		}
		peers := strings.Split(string(content), "\n")
		return c.JSON(fiber.Map{"peerip": peers[:len(peers)-1]})
	})

	fmt.Println("Registry Started")

	log.Fatal(app.Listen(":5050"))
}
