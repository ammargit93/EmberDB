package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Store struct {
	Text map[string]any
}

var (
	StoreStructure = Store{Text: make(map[string]any)}
	mu             sync.Mutex
)

func SaveFile(key string, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := filepath.Dir(path)

		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Println("Error creating directory:", err)
			return
		}

		// Get file name from key
		fileName, _ := GetFile(key)
		fileArr := strings.Split(fileName, "\\")
		fileName = fileArr[len(fileArr)-1]
		log.Println("Saving file:", fileName)

		f, err := os.Create(path)
		if err != nil {
			log.Println("Error creating file:", err)
			return
		}
		if _, err := f.Write([]byte(fileName)); err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
	}
}

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
