package storage

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func InitialiseWAL() {
	_, err := os.Create("../data/wal.log")
	if err != nil {
		fmt.Println(err)
	}
}

var fmu sync.RWMutex

func WriteToWAL(content string) {
	fmu.Lock()
	defer fmu.Unlock()
	file, err := os.OpenFile("../data/wal.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write([]byte(content))
	if err != nil {
		log.Fatalln(err)
	}
}
