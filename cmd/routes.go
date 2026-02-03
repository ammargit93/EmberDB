package main

import (
	"emberdb/internal"
	"emberdb/storage"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Namespace string      `json:"namespace"`
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
}

func stringifyValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func SetKey(c *fiber.Ctx) error {
	var data Response
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	storage.Channel <- "[SETVAL]|" + data.Namespace + "|" + data.Key + "|" + stringifyValue(data.Value) + "|" + internal.InferType(data.Value) + "\n"

	store := &internal.DataStore
	ok, err := store.Insert(data.Namespace, data.Key, data.Value)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !ok {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "key already exists",
		})
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
	md, err := store.Get(namespace, key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"namespace": namespace,
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
	storage.Channel <- "[UPDATEVAL]|" + data.Namespace + "|" + data.Key + "|" + stringifyValue(data.Value) + "|" + internal.InferType(data.Value) + "\n"
	store := &internal.DataStore
	md, err := store.Update(data.Namespace, data.Key, data.Value)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Successfully updated",
		"metadata": md,
	})
}

func DeleteKey(c *fiber.Ctx) error {
	key := c.Params("key")
	namespace := c.Params("namespace")

	storage.Channel <- "[DELETEVAL]|" + namespace + "|" + key + "\n"
	store := &internal.DataStore
	err := store.Delete(namespace, key)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Successfully deleted",
		"key":       key,
		"namespace": namespace,
	})
}

func GetAll(c *fiber.Ctx) error {
	store := &internal.DataStore
	result := store.GetAll()
	return c.JSON(fiber.Map{
		"data": result,
	})
}
