package db

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	Text map[string]any
}

var (
	StoreStructure = Store{Text: make(map[string]any)}
	mu             sync.Mutex
)

func SaveFile(key string, path string) error {
	// Step 1: Get file content from memory (via GetFile)
	fileContentStr, err := GetFile(key)
	if err != nil {
		return fmt.Errorf("error getting file from cache: %w", err)
	}

	// Convert string to byte slice
	fileContent := []byte(fileContentStr)
	mimeType := http.DetectContentType(fileContent)
	var ext string
	switch mimeType {
	case "application/pdf":
		ext = ".pdf"
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "video/mp4":
		ext = ".mp4"
	case "video/x-matroska":
		ext = ".mkv"
	case "video/webm":
		ext = ".webm"
	case "video/ogg":
		ext = ".ogv"
	default:
		ext = ".bin"
	}

	// Add extension if missing
	if filepath.Ext(path) == "" {
		path += ext
	}

	// Step 2: Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	// Step 3: Write the content to the destination file
	if err := os.WriteFile(path, fileContent, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	log.Println("File saved successfully to:", path)
	return nil
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
