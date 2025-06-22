package main

import (
	"emberdb/internal/state"
	"encoding/json"
	"fmt"
	"net/http"
)

type registerResponse struct {
	NodeAddr  string   `json:"node_addr"`
	NodeArray []string `json:"node_array"`
	Leader    string   `json:"leader"`
}

func RegisterNode(port string) error {
	req, err := http.NewRequest("GET", "http://localhost:5050/register", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Port", port)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	state.Mu.Lock()
	state.NodeAddr = response.NodeAddr
	state.AllPeers = response.NodeArray
	state.Leader = response.NodeArray[0]
	state.Mu.Unlock()

	return nil
}
