package internal

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type ValueType uint8

const (
	ValueString ValueType = iota
	ValueInt
	ValueFloat
	ValueBool
	ValueFile // []byte
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
	Type  ValueType `json:"type"`
	Value Value     `json:"value"`
}

type Value struct {
	Type ValueType
	Data []byte
}

func StringValue(v string) Value {
	return Value{Type: ValueString, Data: []byte(v)}
}

func IntValue(v int64) Value {
	return Value{Type: ValueInt, Data: strconv.AppendInt(nil, v, 10)}
}

func FloatValue(v float64) Value {
	return Value{Type: ValueFloat, Data: strconv.AppendFloat(nil, v, 'g', -1, 64)}
}

func BoolValue(v bool) Value {
	if v {
		return Value{Type: ValueBool, Data: []byte{1}}
	}
	return Value{Type: ValueBool, Data: []byte{0}}
}

func FileValue(v []byte) Value {
	return Value{Type: ValueFile, Data: v}
}

func (store *Store) Insert(namespace, key string, value Value) (bool, error) {
	md := Metadata{
		Type:  value.Type,
		Value: value,
	}

	store.Mu.Lock()
	defer store.Mu.Unlock()

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

	if _, exists := nms.Data[key]; exists {
		return false, fiber.NewError(fiber.StatusConflict, "Key exists")
	}

	nms.Data[key] = md
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

func (store *Store) Update(namespace string, key string, value Value) (Metadata, error) {
	md := Metadata{
		Type:  value.Type,
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
