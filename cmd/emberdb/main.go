package main

import (
	"bufio"
	"emberdb/internal/parser"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		msgArr := strings.Split(message, " ")
		output, _ := parser.ParseAndExecute(msgArr)
		conn.Write([]byte(output + "\n"))
	}
}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		panic("Error!!!!")
	}
	fmt.Println("Server running at 8080")
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}
