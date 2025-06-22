package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	NodeStore = []string{}
	mu        sync.Mutex
)

func getClientAddress(r *http.Request) string {
	ip := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = strings.Split(forwarded, ",")[0]
	}
	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		host = ip
	}
	if host == "::1" {
		host = "localhost"
	}
	return host
}

type registerResponse struct {
	NodeAddr  string   `json:"node_addr"`
	NodeArray []string `json:"node_array"`
	Leader    string   `json:"leader"`
}

func main() {
	fmt.Println("Registry server started on :5050")
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientAddress(r)
		port := r.Header.Get("X-Port")
		if port == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "X-Port header is required"})
			return
		}

		nodeAddr := fmt.Sprintf("http://%s%s", clientIP, port)

		mu.Lock()
		for i := 0; i < len(NodeStore); i++ {
			if nodeAddr == NodeStore[i] {
				fmt.Println("Address already included")
				return
			}
		}
		NodeStore = append(NodeStore, nodeAddr)
		leader := NodeStore[0]

		a := NodeStore
		mu.Unlock()

		fmt.Printf("Registered node: %s\n", nodeAddr)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(registerResponse{
			NodeAddr:  nodeAddr,
			NodeArray: a,
			Leader:    leader,
		})
	})

	http.ListenAndServe(":5050", nil)
}
