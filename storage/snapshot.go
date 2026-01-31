package storage

import (
	"emberdb/internal"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const finalPath string = "/data/snapshot.json"

var DurationMap map[string]time.Duration = map[string]time.Duration{
	"s": time.Second,
	"m": time.Second * 60,
	"h": time.Second * 60 * 60,
}

func Snap(duration time.Duration) {
	for {
		if err := saveToJSON(); err != nil {
			fmt.Println("Snapshot failed:", err)
		}
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
	internal.DataStore.Mu.RLock()
	defer internal.DataStore.Mu.RUnlock()

	root, err := projectRoot()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(root, finalPath)
	tmp := fullPath + ".tmp"

	data, err := json.MarshalIndent(internal.DataStore.Namespaces, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	err = os.WriteFile(tmp, data, 0644)
	if err != nil {
		return err
	}
	return os.Rename(tmp, fullPath)
}

func LoadFromJSON() error {
	internal.DataStore.Mu.Lock()
	defer internal.DataStore.Mu.Unlock()

	root, err := projectRoot()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(root, finalPath)

	filebytes, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	var store internal.Store
	if err := json.Unmarshal(filebytes, &store.Namespaces); err != nil {
		return err
	}

	internal.DataStore.Namespaces = store.Namespaces
	return nil
}

func Spawn() {
	internal.DataStore.Mu.RLock()
	defer internal.DataStore.Mu.RUnlock()

	val, exists := internal.ArgMap["snapshot"]

	if !exists {
		fmt.Println("No --snapshot flag, setting fallback interval to 5s")
		val = "5s"
	}
	number := val[:len(val)-1]
	unit := val[len(val)-1:]
	duration, exists := DurationMap[unit]
	if !exists {
		fmt.Println("Invalid time unit")
		return
	}
	drtn, err := strconv.Atoi(number)
	if err != nil {
		fmt.Println(err)
		return
	}

	go Snap(time.Duration(drtn) * duration)

}
