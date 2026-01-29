package storage

import (
	"emberdb/internal"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var SnapshotPath string = "/data/snapshot.json"

func Snap(duration time.Duration) {
	for {
		saveToJSON()
		time.Sleep(duration)
	}
}

var mu sync.RWMutex

func projectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Dir(cwd), nil
}
func saveToJSON() error {
	mu.RLock()
	defer mu.RUnlock()

	root, err := projectRoot()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(root, SnapshotPath)

	data, err := json.MarshalIndent(internal.DataStore, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, data, 0644)
}
