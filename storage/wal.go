package storage

import (
	"fmt"
	"log"
	"os"
)

var Channel chan string = make(chan string)

func InitialiseWAL() error {
	f, err := os.OpenFile("../data/wal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}
func RunWALLoop() {
	for {
		content := <-Channel
		WriteToWAL(content)
		fmt.Println("Written to WAL: ", content)
	}
}

func WriteToWAL(content string) {
	file, err := os.OpenFile("../data/wal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	if err != nil {
		log.Fatalln(err)
	}
	file.Sync()
}
