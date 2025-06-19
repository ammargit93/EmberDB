package main

import (
	"bufio"
	"emberdb/internal/parser"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
)

var (
	mu       sync.Mutex
	NodeAddr string
	AllPeers []string
	Leader   string
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

	fmt.Println("Leader:", response.Leader)
	fmt.Println("Node address:", response.NodeAddr)
	fmt.Println("Node addresses:", response.NodeArray)
	mu.Lock()
	NodeAddr = response.NodeAddr
	AllPeers = response.NodeArray
	Leader = response.NodeArray[0]
	mu.Unlock()

	return nil
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		msgArr := strings.Split(message, " ")
		output, _ := parser.ParseAndExecute(msgArr)
		conn.Write([]byte(output + "\n<END>\n"))
	}
}

func main() {
	port := os.Args[1]
	l, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		panic("Error!!!!")
	}
	fmt.Println("Server running at " + port)
	_ = RegisterNode(port)
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(NodeAddr)
	fmt.Println(AllPeers)
	fmt.Println(Leader)
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}
