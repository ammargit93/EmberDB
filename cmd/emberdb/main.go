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
	"os/signal"
	"strings"
	"syscall"
)

// var (
// 	mu       sync.Mutex
// 	NodeAddr string
// 	AllPeers []string
// 	Leader   string
// )

type peerData struct {
	NodeAddr string   `json:"nodeaddr"`
	AllPeers []string `json:"allpeers"`
	Leader   string   `json:"leader"`
}

type leaderToFollower struct {
	Data   db.Store `json:"data"`
	Sender string   `json:"sender"`
}

func sendRequestToFollowers(data db.Store) error {
	if state.Leader != state.NodeAddr {
		return fmt.Errorf("only the leader should replicate to followers")
	}

	payload := leaderToFollower{
		Data:   data,
		Sender: state.NodeAddr,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Discover peers
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:5050/find_peers", nil)
	req.Header.Set("X-Port", state.NodeAddr[strings.LastIndex(state.NodeAddr, ":"):]) // Extract ":1010" etc
	peerResp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error getting peers:", err)
		return err
	}
	defer peerResp.Body.Close()

	peerBody, err := io.ReadAll(peerResp.Body)
	if err != nil {
		fmt.Println("Error reading peer response:", err)
		return err
	}
	fmt.Println("Peer body ", string(peerBody))
	var peerdata peerData
	if err := json.Unmarshal(peerBody, &peerdata); err != nil {
		fmt.Println("Error unmarshalling peer data:", err)
		return err
	}

	fmt.Println(state.Leader + " Done finding peer leader!")
	fmt.Println(strings.Join(state.AllPeers, " "), "Done finding peer array!")
	fmt.Println(state.NodeAddr + " Done finding peer address!")

	// Send to all followers (excluding self and leader)
	for _, peer := range state.AllPeers {
		if peer == state.NodeAddr || peer == state.Leader {
			continue
		}

		ip := api.SplitIPandIncrementPort(peer)
		resp, err := http.Post(ip+"/replicate", "application/json", strings.NewReader(string(jsonPayload)))
		if err != nil {
			fmt.Println("Error replicating to", peer, ":", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response from", peer, ":", string(body))
		resp.Body.Close()
	}

	return nil
}

func heartbeatSignal() {

	for _, peer := range state.AllPeers {
		if peer == state.NodeAddr || peer == state.Leader {
			continue
		}

		ip := api.SplitIPandIncrementPort(peer)
		resp, err := http.Post(ip+"/heartbeat", "text/plain", strings.NewReader(string("Im alive!")))
		if err != nil {
			fmt.Println("Hearbeat signals to", peer, ":", err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response from", peer, ":", string(body))
		resp.Body.Close()
	}

}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		msgArr := strings.Split(message, " ")
		output, _ := parser.ParseAndExecute(msgArr)
		conn.Write([]byte(output + "\n<END>\n"))
		if strings.ToUpper(msgArr[0]) == "SET" {
			if state.NodeAddr == state.Leader {
				sendRequestToFollowers(db.StoreStructure)
				fmt.Println("Starting replication from leader:", state.NodeAddr)
			} else {
				fmt.Println("Not the leader, skipping replication.")
			}
		}

	}
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down node:", state.NodeAddr)
		if err := DeleteNode(Port); err != nil {
			fmt.Println("Error deleting node from registry:", err)
		} else {
			fmt.Println("Node successfully removed from registry.")
		}
		os.Exit(0)
	}()

	state.Mu.Lock()
	defer state.Mu.Unlock()

	fmt.Println("Leader:", state.Leader)
	fmt.Println("Node address:", state.NodeAddr)
	fmt.Println("Node addresses:", state.AllPeers)
	port := api.IncrementPort(Port)
	go api.StartHTTPServer(port)
	go heartbeatSignal()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}

}
