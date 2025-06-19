package db

import (
	"fmt"
	"os"
	"sync"
)

type Store struct {
	text map[string]any
}

var (
	store = Store{text: make(map[string]any)}
	mu    sync.Mutex
)

func GetFile(key string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	var buf []byte
	buf = store.text[key].([]byte)
	return string(buf), nil
}

func SetFile(key string, value string) (error, bool) {
	filePath := value
	//
	fmt.Println(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error Opening the file", err)
		return err, false
	}
	mu.Lock()
	defer mu.Unlock()
	store.text[key] = data
	return nil, true
}

func GetAllData() string {
	mu.Lock()
	defer mu.Unlock()
	allPairs := ""
	for k, v := range store.text {
		allPairs += k + " : " + v.(string) + "\n"
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
