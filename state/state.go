package state

import "sync"

var (
	Mutx      sync.RWMutex
	DataStore = make(map[string]any)
)
