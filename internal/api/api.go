package api

import (
	"emberdb/internal/db"
	"emberdb/internal/state"
	"encoding/json"
	"fmt"
	"net/http"
)

type followerResponse struct {
	Data   db.Store `json:"data"`
	Sender string   `json:"sender"`
}

type Data struct {
	DataMap db.Store `json:"store"`
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
		for i := 0; i < len(state.AllPeers); i++ {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(&db.StoreStructure)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		// w.Header().Set("Content-Type", "application/json")
		// err := json.NewEncoder(w).Encode(&data)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }

	})
	http.ListenAndServe(port, nil)
}
