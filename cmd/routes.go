package main

import (
	"emberdb/internal"
	"emberdb/storage"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Namespace string      `json:"namespace"`
	Key       string      `json:"key"`
	Type      string      `json:"type"`
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
	val, err := internal.ParseValue(data.Type, data.Value)

	storage.Channel <- "[SETVAL]|" +
		data.Namespace + "|" +
		data.Key + "|" +
		internal.StringifyValue(val) + "\n"

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	ok, err := store.Insert(data.Namespace, data.Key, val)

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
		"value":     string(md.Value.Data),
	})
}

func UpdateKey(c *fiber.Ctx) error {
	var data Response
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	val, err := internal.ParseValue(data.Type, data.Value)
	store := &internal.DataStore
	md, err := store.Update(data.Namespace, data.Key, val)

	storage.Channel <- "[UPDATEVAL]|" +
		data.Namespace + "|" +
		data.Key + "|" +
		internal.StringifyValue(val) + "\n"

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
