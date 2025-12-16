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
		mutx.Lock()
		ActivePeers = parsed.PeerIP
		mutx.Unlock()

		fmt.Println("Active peers:", ActivePeers)

		time.Sleep(5 * time.Second)
	}
}

type Data struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

var DataStore = make(map[string]any)

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

	app.Post("/set", func(c *fiber.Ctx) error {
		var data Data
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}
		mutx.Lock()
		defer mutx.Unlock()
		if _, exists := DataStore[data.Key]; exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "key already exists",
				"key":   data.Key,
			})
		}
		DataStore[data.Key] = data.Value
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"key":   data.Key,
			"value": data.Value,
		})
	})

	app.Get("/get/:key", func(c *fiber.Ctx) error {
		mutx.Lock()
		val := DataStore[c.Params("key")]
		mutx.Unlock()
		return c.JSON(fiber.Map{
			"value": val,
		})
	})

	app.Patch("/update", func(c *fiber.Ctx) error {
		var data Data
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}
		key := data.Key
		value := data.Value
		mutx.Lock()
		defer mutx.Unlock()
		if _, exists := DataStore[key]; !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "key does not exist",
				"key":   key,
			})
		}
		DataStore[key] = value
		return c.JSON(fiber.Map{
			"message": "Successfully updated",
			"key":     key,
			"value":   value,
		})
	})

	app.Delete("/delete/:key", func(c *fiber.Ctx) error {
		key := c.Params("key")
		mutx.Lock()
		defer mutx.Unlock()
		if _, exists := DataStore[key]; !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "key does not exist",
				"key":   key,
			})
		}
		delete(DataStore, key)
		return c.JSON(fiber.Map{
			"message": "Successfully deleted",
			"key":     key,
		})
	})

	app.Get("/getall", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"Data": DataStore,
		})
	})

	fmt.Println("Server :" + port + " started")

	app.Listen(":" + port)
}
