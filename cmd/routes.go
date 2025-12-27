package main

import "github.com/gofiber/fiber/v2"

func SetKey(c *fiber.Ctx) error {
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
}

func GetKey(c *fiber.Ctx) error {
	mutx.Lock()
	val := DataStore[c.Params("key")]
	mutx.Unlock()
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
}

func DeleteKey(c *fiber.Ctx) error {
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
}

func GetAll(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Data": DataStore,
	})
}

func MSet(c *fiber.Ctx) error {
	var data map[string]any
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	mutx.Lock()
	defer mutx.Unlock()
	for k, v := range data {
		DataStore[k] = v
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
	mutx.RLock()
	defer mutx.RUnlock()
	res := make(map[string]any)
	for _, k := range keys {
		if v, ok := DataStore[k]; ok {
			res[k] = v
		} else {
			res[k] = nil
		}
	}
	return c.JSON(fiber.Map{
		"values": res,
	})
}
