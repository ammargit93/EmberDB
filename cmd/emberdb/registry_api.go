package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// type registerResponse struct {
// 	NodeAddr  string   `json:"node_addr"`
// 	NodeArray []string `json:"node_array"`
// 	Leader    string   `json:"leader"`
// }

// func RegisterNode(port string) error {
// 	req, err := http.NewRequest("GET", "http://localhost:5050/register", nil)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("X-Port", port)

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	var response registerResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
// 		return fmt.Errorf("failed to decode response: %v", err)
// 	}

// 	fmt.Println("Leader:", response.Leader)
// 	fmt.Println("Node address:", response.NodeAddr)
// 	fmt.Println("Node addresses:", response.NodeArray)
// 	if response.Leader == response.NodeAddr {
// 		mu.Lock()
// 		NodeAddr = response.NodeAddr
// 		AllPeers = response.NodeArray
// 		Leader = response.Leader
// 		mu.Unlock()
// 	}
// 	return nil
// }
