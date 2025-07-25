package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func callConnectAPI(command string) (string, error) {
	url := "http://localhost:8001/connect-db"
	resp, err := http.Post(url, "text/plain", strings.NewReader(command))
	if err != nil {
		log.Println("Failed to call connect API:", err)
		return "Error" + err.Error(), err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return "Error" + err.Error(), err
	}
	return string(body), err
}

func main() {
	start := time.Now()

	output, _ := callConnectAPI("SAVE file video.mp4")

	elapsed := time.Since(start)
	fmt.Println(output)
	fmt.Printf("⏱️ Time taken: %.2f seconds\n", elapsed.Seconds())
}
