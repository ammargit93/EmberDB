package sdk

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// var conn net.Conn
type Client struct {
	conn net.Conn
}

func Connect(addr string) Client {
	connection, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Connection failed:", err)
		os.Exit(1)
	}
	var client Client
	client.conn = connection
	return client
}

func (client Client) SetValue(key string, value any) error {
	valString, ok := value.(string)

	if !ok {
		fmt.Println("Value is not a string")
	}

	command := "SET " + key + " " + valString

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}

	return nil

}

func (client Client) GetValue(key string) (any, error) {

	command := "GET " + key

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return -1, err
	}
	reader := bufio.NewReader(client.conn)
	line, err := reader.ReadString('\n')

	return line, err

}

func (client Client) DelValue(key string) error {

	command := "DEL " + key

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}

	return err
}

func (client Client) UpdateValue(key string, value string) error {

	command := "UPDATE " + key + " " + value

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}

	return err
}

func (client Client) SetFile(key string, value string) error {

	command := "SETFILE " + key + " " + value

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}

	return err
}

func (client Client) GetFile(key string) (string, error) {

	command := "GETFILE " + key

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err.Error(), err
	}
	reader := bufio.NewReader(client.conn)
	line, err := reader.ReadString('\n')

	return line, err
}
