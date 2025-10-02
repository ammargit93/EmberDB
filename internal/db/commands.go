package db

import (
	"fmt"
	"os"
	"sync"
)

type Store struct {
	Text map[string]any
}

var (
	StoreStructure = Store{Text: make(map[string]any)}
	mu             sync.Mutex
)

func GetFile(key string) (string, error) {
	mu.Lock()
	defer mu.Unlock()
	var buf []byte
	buf = StoreStructure.Text[key].([]byte)
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
	StoreStructure.Text[key] = data
	return nil, true
}

func SaveFile(key string, fileName string) error {
	fileContent, _ := GetFile(key)

	fmt.Println(fileContent)
	os.WriteFile(fileName, []byte(fileContent), 0755)
	return nil
}

func GetAllData() string {
	mu.Lock()
	defer mu.Unlock()
	allPairs := ""
	for k, v := range StoreStructure.Text {
		allPairs += k + " : " + v.(string) + "\n"
	}
	// fmt.Println(allPairs)
	return allPairs
}

func UpdateValue(key string, value string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := StoreStructure.Text[key]
	if exists {
		StoreStructure.Text[key] = value
		return true
	}
	return false
}

func SetValue(key string, value string) bool {
	mu.Lock()
	defer mu.Unlock()
	_, exists := StoreStructure.Text[key].(string)
	if exists {
		return true
	} else {
		StoreStructure.Text[key] = value
		return false
	}

}

func GetValue(key string) any {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range StoreStructure.Text {
		if key == k {
			return v
		}
	}
	return "No such key exists in the store."

}

func DeleteKey(key string) bool {
	mu.Lock()
	defer mu.Unlock()
	for k, _ := range StoreStructure.Text {
		if k == key {
			delete(StoreStructure.Text, key)
			return true
		}
	}
	return false
}
