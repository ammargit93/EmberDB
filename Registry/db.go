package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "./Registry/registry.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS peers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT UNIQUE NOT NULL
	);`)
	if err != nil {
		log.Fatal(err)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientAddress(r)
	port := r.Header.Get("X-Port")
	if port == "" {
		http.Error(w, "X-Port header required", http.StatusBadRequest)
		return
	}
	nodeAddr := fmt.Sprintf("http://%s%s", clientIP, port)

	_, err := db.Exec(`INSERT OR IGNORE INTO peers(address) VALUES (?)`, nodeAddr)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`SELECT address FROM peers`)
	if err != nil {
		http.Error(w, "DB read error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var peers []string
	for rows.Next() {
		var addr string
		rows.Scan(&addr)
		peers = append(peers, addr)
	}

	leader := peers[0] // First one is leader (simplified logic)
	resp := registerResponse{
		NodeAddr:  nodeAddr,
		NodeArray: peers,
		Leader:    leader,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func FindAllPeers(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientAddress(r)
	port := r.Header.Get("X-Port")
	if port == "" {
		http.Error(w, "X-Port header required", http.StatusBadRequest)
		return
	}
	nodeAddr := fmt.Sprintf("http://%s%s", clientIP, port)
	rows, err := db.Query(`SELECT address FROM peers`)
	if err != nil {
		http.Error(w, "DB read error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var peers []string
	for rows.Next() {
		var addr string
		rows.Scan(&addr)
		peers = append(peers, addr)
	}

	leader := peers[0] // First one is leader (simplified logic)
	resp := registerResponse{
		NodeAddr:  nodeAddr,
		NodeArray: peers,
		Leader:    leader,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	// json.NewEncoder(w).Encode(resp)
}
