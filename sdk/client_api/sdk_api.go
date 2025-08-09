package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func connect(command string) string {
	address := "localhost:1010"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}
	defer conn.Close()
	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return ""
	}
	reader := bufio.NewReader(conn)
	lines := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "<END>" {
			break
		}
		// fmt.Print(line)
		lines += line
	}
	return lines
}

func ConnectToLeader() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyRead, _ := io.ReadAll(r.Body)
		command := string(bodyRead)
		log.Println(command)
		conn := connect(command)
		w.Write([]byte(conn))
	}

}

func main() {
	http.HandleFunc("/connect-db", ConnectToLeader())
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Println("Failed to start HTTP server:", err)
	}
}
