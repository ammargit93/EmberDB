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

	http.HandleFunc("/find_leader", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		response := map[string]interface{}{
			"leader": "",
			"error":  nil,
		}

		status := http.StatusOK
		if len(NodeStore) == 0 {
			status = http.StatusNotFound
			response["error"] = "no leader available"
		} else {
			response["leader"] = NodeStore[0]
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":5050", nil)
}
