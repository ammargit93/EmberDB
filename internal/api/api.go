package api

import (
	"emberdb/internal/db"
	"encoding/json"
	"fmt"
	"net/http"
)

type followerResponse struct {
	Data   db.Store `json:"data"`
	Sender string   `json:"sender"`
}

func StartHTTPServer(port string) {
	http.HandleFunc("/replicate", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var fr followerResponse
		if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			fmt.Println("JSON decode error:", err)
			return
		}

		db.StoreStructure = fr.Data
		fmt.Println("Replication request received from:", fr.Sender)
		fmt.Println("Updated store:", db.StoreStructure)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "replicated"})
	})

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
