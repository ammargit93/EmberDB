package main

import (
	"emberdb/internal/state"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type registerResponse struct {
	NodeAddr  string   `json:"nodeaddr"`
	NodeArray []string `json:"nodearray"`
	Leader    string   `json:"leader"`
}

func RegisterNode(port string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:5050/register", nil)
	req.Header.Set("X-Port", port)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Registry server error:", string(body))
		return fmt.Errorf("register failed with status: %d", resp.StatusCode)
	}

	var regResp registerResponse
	if err := json.NewDecoder(resp.Body).Decode(&regResp); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	state.Mu.Lock()
	state.NodeAddr = regResp.NodeAddr
	state.AllPeers = regResp.NodeArray
	state.Leader = regResp.Leader
	state.Mu.Unlock()

	return nil
}
