package main

import (
	"bufio"
	"emberdb/internal/db"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		msgArr := strings.Split(message, " ")
		op := msgArr[0]

		if op == "SET" {
			db.SetValue(msgArr[1], msgArr[2])
			conn.Write([]byte("SET OK\n"))
		} else if op == "GET" && len(msgArr) == 2 {
			value := db.GetValue(msgArr[1])
			conn.Write([]byte(value.(string) + "\n"))
		} else if op == "DEL" && len(msgArr) == 2 {
			value := db.DeleteKey(msgArr[1])
			conn.Write([]byte(value + "\n"))
		}

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
