package storage

import (
	"bufio"
	"emberdb/internal"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

var mapper map[string]string = map[string]string{
	"0": "string",
	"1": "int",
	"2": "float",
	"3": "bool",
	"4": "file",
}

func ReplayWAL() error {
	f, err := os.Open("../data/wal.log")
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	store := &internal.DataStore

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue // corrupted line
		}

		op := parts[0]
		ns := parts[1]
		key := parts[2]

		switch op {

		case "[SETVAL]", "[UPDATEVAL]":
			if len(parts) != 4 {
				continue
			}

			tv := strings.SplitN(parts[3], ":", 2)
			if len(tv) != 2 {
				continue
			}

			typ, err := strconv.Atoi(tv[0])
			if err != nil {
				continue
			}
			data, err := base64.StdEncoding.DecodeString(tv[1])
			if err != nil {
				continue
			}
			val := internal.Value{
				Type: internal.ValueType(typ),
				Data: data,
			}

			if op == "[SETVAL]" {
				store.Insert(ns, key, val)
			} else {
				store.Update(ns, key, val)
			}

		case "[DELETE]":
			store.Delete(ns, key)
		}
	}

	return scanner.Err()
}
