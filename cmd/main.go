package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

func register(port string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:5050/register", nil)
	req.Header.Set("X-Port", port)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}

var ActivePeers []string
var mutx sync.RWMutex

type FetchPeersResponse struct {
	PeerIP []string `json:"peerip"`
}

func fetchPeers() {
	for {
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "http://localhost:5050/fetch-peers", nil)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Println("Read error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		var parsed FetchPeersResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			fmt.Println("JSON error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// update global slice
		mutx.Lock()
		ActivePeers = parsed.PeerIP
		mutx.Unlock()

		fmt.Println("Active peers:", ActivePeers)

		time.Sleep(5 * time.Second)
	}
}

func main() {
	port := os.Args[1]
	app := fiber.New()
	register(port)
	go fetchPeers()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from " + port,
		})
	})

	fmt.Println("Server :" + port + " started")
	app.Listen(":" + port)
}
