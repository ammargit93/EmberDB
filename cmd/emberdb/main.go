package main

import (
	"bufio"
	"emberdb/internal/api"
	"emberdb/internal/db"
	"emberdb/internal/parser"
	"emberdb/internal/state"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// var (
// 	mu       sync.Mutex
// 	NodeAddr string
// 	AllPeers []string
// 	Leader   string
// )

type followerToLeader struct {
	Data   db.Store `json:"data"`
	Sender string   `json:"sender"`
}

func splitIPandIncrementPort(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) < 3 {
		return "http://localhost:6060"
	}
	host := parts[1]    // "//localhost"
	portStr := parts[2] // "9090"
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return "http://localhost:6060"
	}
	newPort := portInt + 1
	host = strings.TrimPrefix(host, "//")
	return fmt.Sprintf("http://%s:%d", host, newPort)
}

func sendRequestToLeader(data db.Store) error {
	if state.Leader == state.NodeAddr {
		return fmt.Errorf("this node is the leader; no need to forward")
	}
	payload := followerToLeader{
		Data:   data,
		Sender: state.NodeAddr,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	peerResp, err := http.Get("http://localhost:5050/find_peers")
	if err != nil {
		fmt.Println("Error is peerResp", err)
	}
	peerbdy, err := io.ReadAll(peerResp.Body)
	if err != nil {
		fmt.Println("Error is ", err)
	}
	fmt.Println(string(peerbdy) + "Done finding peers!")
	ip := splitIPandIncrementPort(state.Leader)
	resp, err := http.Post(ip+"/replicate", "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		fmt.Println("Response Error is ", err)
	}
	bdy, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error is ", err)
	}
	fmt.Println(string(bdy) + "Done forwarding!")

	if err != nil {
		return err
	}
	defer peerResp.Body.Close()
	defer resp.Body.Close()
	return nil
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		msgArr := strings.Split(message, " ")
		output, _ := parser.ParseAndExecute(msgArr)
		conn.Write([]byte(output + "\n<END>\n"))
		sendRequestToLeader(db.StoreStructure)
	}
}
func IncrementPort(portStr string) string {
	portNum := strings.TrimPrefix(portStr, ":")
	portInt, err := strconv.Atoi(portNum)
	if err != nil {
		return ":6060"
	}

	return fmt.Sprintf(":%d", portInt+1)
}

func main() {
	Port := os.Args[1]
	l, err := net.Listen("tcp", Port)

	if err != nil {
		fmt.Println(err)
		panic("Error!!!!")
	}
	fmt.Println("Server running at " + Port)

	_ = RegisterNode(Port)
	state.Mu.Lock()
	defer state.Mu.Unlock()

	fmt.Println("Leader:", state.Leader)
	fmt.Println("Node address:", state.NodeAddr)
	fmt.Println("Node addresses:", state.AllPeers)
	port := IncrementPort(Port)
	go api.StartHTTPServer(port)

	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}
