package db

import (
	"fmt"
	"sync"
)

type Store struct {
	text map[string]any
}

var (
	store = Store{text: make(map[string]any)}
	mu    sync.Mutex
)

func SetValue(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	_, exists := store.text[key].(string)
	if exists {
		fmt.Println("Key already exists")
		return
	}
	store.text[key] = value
}

func GetValue(key string) any {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range store.text {
		if key == k {
			return v
		}
	}
	return "No such key exists in the store."

}

func DeleteKey(key string) string {
	mu.Lock()
	defer mu.Unlock()
	for k, _ := range store.text {
		if k == key {
			delete(store.text, key)
			return "Key Deleted Successfully"
		}
	}
	return "Key could not be deleted"
}
