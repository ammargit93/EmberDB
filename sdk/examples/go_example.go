package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func callConnectAPI(command string) string {
	url := "http://localhost:8001/connect-db"
	resp, err := http.Post(url, "text/plain", strings.NewReader(command))
	if err != nil {
		log.Println("Failed to call connect API:", err)
		return "Error" + err.Error()
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return "Error" + err.Error()
	}
	// fmt.Println("Response from emberdb: ")
	return string(body)
}
func main() {
	output := callConnectAPI("SET d 100") // output: SET OK
	fmt.Println(output)
	output = callConnectAPI("GET d") // output: 100
	fmt.Println(output)
}
