package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

func flushBuffer(client Client) error {
	reader := bufio.NewReader(client.conn)
	_, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}
	return nil
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
	_ = flushBuffer(client)

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
	line = strings.Replace(line, "\n", "", -1)

	return line, err

}

func (client Client) DelValue(key string) error {

	command := "DEL " + key

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}
	_ = flushBuffer(client)

	return err
}

func (client Client) UpdateValue(key string, value string) error {

	command := "UPDATE " + key + " " + value

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}
	_ = flushBuffer(client)

	return err
}

func (client Client) SetFile(key string, value string) error {

	command := "SETFILE " + key + " " + value

	_, err := client.conn.Write([]byte(command + "\n"))
	if err != nil {
		fmt.Println("Error sending command:", err)
		return err
	}
	_ = flushBuffer(client)

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
	lineArr := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.TrimSpace(line) == "<END>" {
			break
		}
		lineArr = append(lineArr, line)
	}

	return strings.Join(lineArr, "\n"), err
}
