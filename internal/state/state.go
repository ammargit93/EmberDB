package state

import "sync"

var (
	Mu       sync.Mutex
	NodeAddr string
	AllPeers []string
	Leader   string
)
