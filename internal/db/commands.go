package db

import (
	"sync"
)

type Store struct {
	text map[string]any
}

var (
	store = Store{text: make(map[string]any)}
	mu    sync.Mutex
)

func GetAllData() string {
	mu.Lock()
	defer mu.Unlock()
	allPairs := ""
	for k, v := range store.text {
		allPairs += k + " : " + v.(string) + "\\"
	}
	// fmt.Println(allPairs)
	return allPairs
}

func UpdateValue(key string, value string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := store.text[key]
	if exists {
		store.text[key] = value
		return true
	}
	return false
}

func SetValue(key string, value string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := store.text[key].(string)
	if exists {
		return true
	} else {
		store.text[key] = value
		return false
	}

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

func DeleteKey(key string) bool {
	mu.Lock()
	defer mu.Unlock()
	for k, _ := range store.text {
		if k == key {
			delete(store.text, key)
			return true
		}
	}
	return false
}
