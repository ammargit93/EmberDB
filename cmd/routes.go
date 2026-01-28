package main

import (
	"emberdb/internal"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var mu sync.RWMutex

type Response struct {
	Namespace string      `json:"namespace"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

func SetKey(c *fiber.Ctx) error {
	var data Response
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	store := &internal.DataStore
	ok, err := store.Insert(data.Namespace, data.Key, data.Value)
	if !ok {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"namespace": data.Namespace,
		"key":       data.Key,
		"value":     data.Value,
	})
}

func GetKey(c *fiber.Ctx) error {

	key := c.Params("key")
	namespace := c.Params("namespace")

	store := &internal.DataStore
	md, _ := store.Get(namespace, key)
	return c.JSON(fiber.Map{
		"namespace": md,
		"key":       key,
		"value":     md.Value,
	})
}

func UpdateKey(c *fiber.Ctx) error {
	var data Response
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	store := &internal.DataStore
	md, _ := store.Update(data.Namespace, data.Key, data.Value)
	return c.JSON(fiber.Map{
		"message":   "Successfully updated",
		"namespace": store.Namespaces[data.Namespace],
		"metadata":  md,
	})
}

func DeleteKey(c *fiber.Ctx) error {
	key := c.Params("key")
	namespace := c.Params("namespace")
	store := &internal.DataStore
	store.Delete(namespace, key)
	return c.JSON(fiber.Map{
		"message":   "Successfully deleted",
		"key":       key,
		"namespace": store.Namespaces[namespace],
	})
}

func GetAll(c *fiber.Ctx) error {
	store := &internal.DataStore
	result := store.GetAll()
	return c.JSON(fiber.Map{
		"Data": result,
	})
}
