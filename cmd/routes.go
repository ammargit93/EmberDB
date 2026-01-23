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
	// create Metadata
	md := internal.Metadata{
		Type:  internal.Datatype(internal.InferType(data.Value)),
		Value: data.Value,
	}
	// create namespace if not exists
	Namespace := internal.Namespace{
		Name: data.Namespace,
		Data: map[string]internal.Metadata{
			data.Key: md,
		},
	}
	mu.Lock()

	store := &internal.DataStore
	if store.Namespaces == nil {
		store.Namespaces = make(map[string]*internal.Namespace)
	}
	store.Namespaces[data.Namespace] = &Namespace

	defer mu.Unlock()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"namespace": data.Namespace,
		"key":       data.Key,
		"value":     data.Value,
	})
}

func GetKey(c *fiber.Ctx) error {

	key := c.Params("key")
	namespace := c.Params("namespace")

	mu.RLock()
	nms, ok := internal.DataStore.Namespaces[namespace]
	if !ok {
		return fiber.NewError(fiber.StatusNotFound, "Namespace not found")
	}
	md, ok := nms.Data[key]
	if !ok {
		return fiber.NewError(fiber.StatusNotFound, "key not found")
	}
	mu.RUnlock()

	return c.JSON(fiber.Map{
		"namespace": nms,
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
	namespace := data.Namespace
	mu.Lock()
	defer mu.Unlock()

	// create Metadata
	md := internal.Metadata{
		Type:  internal.Datatype(internal.InferType(data.Value)),
		Value: data.Value,
	}
	// create namespace if not exists
	Namespace := internal.Namespace{
		Name: namespace,
		Data: map[string]internal.Metadata{
			data.Key: md,
		},
	}

	store := &internal.DataStore
	if store.Namespaces == nil {
		store.Namespaces = make(map[string]*internal.Namespace)
		return fiber.NewError(fiber.StatusNotFound, "Cannot update uninitialised store.")
	}
	store.Namespaces[namespace] = &Namespace

	return c.JSON(fiber.Map{
		"message":   "Successfully updated",
		"namespace": store.Namespaces[namespace],
	})
}

// func DeleteKey(c *fiber.Ctx) error {
// 	key := c.Params("key")
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if _, exists := state.DataStore[key]; !exists {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": "key does not exist",
// 			"key":   key,
// 		})
// 	}
// 	delete(state.DataStore, key)
// 	return c.JSON(fiber.Map{
// 		"message": "Successfully deleted",
// 		"key":     key,
// 	})
// }

// func GetAll(c *fiber.Ctx) error {
// 	fmt.Println(state.DataStore)
// 	return c.JSON(fiber.Map{
// 		"Data": state.DataStore,
// 	})
// }

// func MSet(c *fiber.Ctx) error {
// 	var data map[string]any
// 	if err := c.BodyParser(&data); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "invalid request body",
// 		})
// 	}
// 	mu.Lock()
// 	defer mu.Unlock()
// 	for k, v := range data {
// 		state.DataStore[k] = v
// 	}
// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"set": data,
// 	})
// }

// func MGet(c *fiber.Ctx) error {
// 	var keys []string
// 	if err := c.BodyParser(&keys); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "invalid request body, expected JSON array of keys",
// 		})
// 	}
// 	mu.RLock()
// 	defer mu.RUnlock()
// 	res := make(map[string]any)
// 	for _, k := range keys {
// 		if v, ok := state.DataStore[k]; ok {
// 			res[k] = v
// 		} else {
// 			res[k] = nil
// 		}
// 	}
// 	return c.JSON(fiber.Map{
// 		"values": res,
// 	})
// }
