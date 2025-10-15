package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func connect(command string) {
	port := os.Args[1]
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return
	}

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "<END>" {
			break
		}
		fmt.Print(line)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("ember> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "quit" || line == "q" {
			break
		}
		connect(line)
	}
}
