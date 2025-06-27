package main

import (
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

type registerResponse struct {
	NodeAddr  string   `json:"nodeaddr"`
	NodeArray []string `json:"nodearray"`
	Leader    string   `json:"leader"`
}

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

func main() {
	InitDB()
	fmt.Println("Registry server started on :5050")
	http.HandleFunc("/register", Register)
	http.HandleFunc("/find_peers", FindAllPeers)
	http.ListenAndServe(":5050", nil)
}
