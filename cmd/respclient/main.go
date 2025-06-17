package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func connect(command string) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	conn.Write([]byte(command + "\n"))

	response, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		fmt.Println("Connection Error: ", err)
		return
	}
	fmt.Print(response)
	defer conn.Close()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("ember> ")
		line, _ := reader.ReadString(byte('\n'))
		if line == "quit" || line == "q" {
			os.Exit(1)
			break
		}
		connect(line)
	}
}
