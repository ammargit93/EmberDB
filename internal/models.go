package internal

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Datatype string

const (
	TypeString Datatype = "string"
	TypeInt    Datatype = "int"
	TypeFloat  Datatype = "float"
	TypeBool   Datatype = "bool"
	TypeFile   Datatype = "file"
)

type Store struct {
	Mu         sync.RWMutex
	Namespaces map[string]*Namespace `json:"namespaces"` // map for fast lookup
}

type Namespace struct {
	Name string              `json:"name"`
	Data map[string]Metadata `json:"data"` // key -> Metadata
}

type Metadata struct {
	Type  Datatype    `json:"type"`
	Value interface{} `json:"value"`
}

func InferType(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case float64:
		return "float"
	case float32:
		return "float"
	case string:
		return "string"
	case bool:
		return "bool"
	default:
		return "file"
	}
}

var mu sync.RWMutex

func (store *Store) Insert(namespace string, key string, value interface{}) (bool, error) {
	md := Metadata{
		Type:  Datatype(InferType(value)),
		Value: value,
	}
	store.Mu.Lock()
	defer store.Mu.Unlock()

	// create namespace if not exists
	if store.Namespaces == nil {
		store.Namespaces = make(map[string]*Namespace)
	}

	nms, exists := store.Namespaces[namespace]
	if !exists {
		nms = &Namespace{
			Name: namespace,
			Data: make(map[string]Metadata),
		}
		store.Namespaces[namespace] = nms

	}
	_, exists = nms.Data[key]
	if !exists {
		nms.Data[key] = md
	} else {
		return false, fiber.NewError(fiber.StatusConflict, "Key exists")
	}
	return true, nil
}

func (store *Store) Get(namespace string, key string) (Metadata, error) {
	store.Mu.RLock()
	defer store.Mu.RUnlock()

	nms, exists := store.Namespaces[namespace]
	if !exists {
		return Metadata{}, fiber.NewError(fiber.StatusNotFound, "Namespace not found")
	}

	md, exists := nms.Data[key]
	if !exists {
		return Metadata{}, fiber.NewError(fiber.StatusNotFound, "Key not found")
	}

	return md, nil
}

func (store *Store) Update(namespace string, key string, value interface{}) (Metadata, error) {
	md := Metadata{
		Type:  Datatype(InferType(value)),
		Value: value,
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

	if store.Namespaces == nil {
		return Metadata{}, fiber.NewError(fiber.StatusNotFound, "Store uninitialized")
	}

	nms, exists := store.Namespaces[namespace]
	if !exists {
		return Metadata{}, fiber.NewError(fiber.StatusNotFound, "Namespace not found")
	}

	nms.Data[key] = md
	return md, nil
}

func (store *Store) Delete(namespace string, key string) error {
	store.Mu.Lock()
	defer store.Mu.Unlock()

	nms, exists := store.Namespaces[namespace]
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, "Namespace not found")
	}

	if _, exists := nms.Data[key]; !exists {
		return fiber.NewError(fiber.StatusNotFound, "Key not found")
	}

	delete(nms.Data, key)
	return nil
}

func (store *Store) GetAll() map[string]*Namespace {
	store.Mu.RLock()
	defer store.Mu.RUnlock()

	if store.Namespaces == nil {
		return make(map[string]*Namespace)
	}

	copyMap := make(map[string]*Namespace, len(store.Namespaces))
	for ns, n := range store.Namespaces {
		newNS := &Namespace{
			Name: n.Name,
			Data: make(map[string]Metadata, len(n.Data)),
		}
		for k, v := range n.Data {
			newNS.Data[k] = v
		}
		copyMap[ns] = newNS
	}

	return copyMap
}

// db
// └── Namespaces (map[string]*Namespace)
//     ├── "users" → *Namespace
//     │       Name: "users"
//     │       Data:
//     │           "username" → Metadata{Type: "string", Value: "ammar"}
//     │           "age"      → Metadata{Type: "int", Value: 23}
//     └── "products" → *Namespace
//             Name: "products"
//             Data:
//                 "laptop" → Metadata{Type: "string", Value: "MacBook"}
//                 "price"  → Metadata{Type: "float", Value: 2999.99}
