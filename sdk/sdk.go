package sdk

import (
	"fmt"
	"net"
	"os"
)

var conn net.Conn

func Connect(addr string) {
	connection, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}
	conn = connection
}

func SetValue(key string, value any) error {
	valString, ok := value.(string)

	if !ok {
		fmt.Println("Value is not a string")
	}

	command := "SET " + key + " " + valString

	_, err := conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}

	return nil

}
