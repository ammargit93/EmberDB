package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

type Store struct {
	text []map[string]string
}

var (
	store Store
	mu    sync.Mutex
)

func setValue(key string, value string) {
	mu.Lock()
	defer mu.Unlock()
	mapElement := map[string]string{key: value}
	store.text = append(store.text, mapElement)
}

func getValue(key string) string {
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < len(store.text); i++ {
		m := store.text[i]
		value, exists := m[key]
		if exists {
			return value
		}
	}
	fmt.Println("No such key exists in the store.")
	return ""
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		msgArr := strings.Split(message, " ")
		op := msgArr[0]
		if op == "SET" {
			setValue(msgArr[1], msgArr[2])
			fmt.Println("SET OK")
		} else if op == "GET" && len(msgArr) == 2 {
			value := getValue(msgArr[1])
			fmt.Println(value)
		}
		// fmt.Println(store.text)

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
