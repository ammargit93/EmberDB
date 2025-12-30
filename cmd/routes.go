package main

import (
	"emberdb/state"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var mu sync.RWMutex

func SetKey(c *fiber.Ctx) error {
	var data Data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	mu.Lock()
	defer mu.Unlock()
	if _, exists := state.DataStore[data.Key]; exists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "key already exists",
			"key":   data.Key,
		})
	}
	state.DataStore[data.Key] = data.Value
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"key":   data.Key,
		"value": data.Value,
	})
}

func GetKey(c *fiber.Ctx) error {
	mu.Lock()
	val := state.DataStore[c.Params("key")]
	mu.Unlock()
	return c.JSON(fiber.Map{
		"value": val,
	})
}

func UpdateKey(c *fiber.Ctx) error {
	var data Data
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	key := data.Key
	value := data.Value
	mu.Lock()
	defer mu.Unlock()
	if _, exists := state.DataStore[key]; !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "key does not exist",
			"key":   key,
		})
	}
	state.DataStore[key] = value
	return c.JSON(fiber.Map{
		"message": "Successfully updated",
		"key":     key,
		"value":   value,
	})
}

func DeleteKey(c *fiber.Ctx) error {
	key := c.Params("key")
	mu.Lock()
	defer mu.Unlock()
	if _, exists := state.DataStore[key]; !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "key does not exist",
			"key":   key,
		})
	}
	delete(state.DataStore, key)
	return c.JSON(fiber.Map{
		"message": "Successfully deleted",
		"key":     key,
	})
}

func GetAll(c *fiber.Ctx) error {
	fmt.Println(state.DataStore)
	return c.JSON(fiber.Map{
		"Data": state.DataStore,
	})
}

func MSet(c *fiber.Ctx) error {
	var data map[string]any
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	mu.Lock()
	defer mu.Unlock()
	for k, v := range data {
		state.DataStore[k] = v
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"set": data,
	})
}

func MGet(c *fiber.Ctx) error {
	var keys []string
	if err := c.BodyParser(&keys); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body, expected JSON array of keys",
		})
	}
	mu.RLock()
	defer mu.RUnlock()
	res := make(map[string]any)
	for _, k := range keys {
		if v, ok := state.DataStore[k]; ok {
			res[k] = v
		} else {
			res[k] = nil
		}
	}
	return c.JSON(fiber.Map{
		"values": res,
	})
}
