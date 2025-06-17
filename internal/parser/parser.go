package parser

import (
	"emberdb/internal/db"
	"strings"
)

var (
	commandParser = map[string]CommandFunc{
		"SET":    handleSet,
		"GET":    handleGet,
		"DEL":    handleDelete,
		"UPDATE": handleUpdate,
		"GETALL": getAllPairs,
	}
)

type CommandFunc func(args []string) string

func getAllPairs(args []string) string {
	data := db.GetAllData()
	// s := strings.Replace(data, "\\", "\n", len(data))
	return data
}

func handleSet(args []string) string {
	key := args[0]
	if len(args) > 2 {
		value := strings.Join(args[1:], " ")
		if db.SetValue(key, value) {
			return "Key already exists"
		} else {
			return "SET OK"
		}
	} else {
		value := args[1]
		if db.SetValue(key, value) {
			return "Key already exists"
		} else {
			return "SET OK"
		}
	}

	// return "SET OK"
}

func handleGet(args []string) string {
	key := args[0]
	value := db.GetValue(key).(string)
	return value
}

func handleDelete(args []string) string {
	key := args[0]
	isDeleted := db.DeleteKey(key)
	if isDeleted {
		return "Key Deleted Successfully"
	} else {
		return "Key could not be deleted"
	}
}

func handleUpdate(args []string) string {
	key := args[0]
	value := args[1]
	isUpdated := db.UpdateValue(key, value)
	if isUpdated {
		return "Key Updated Successfully"
	} else {
		return "Key could not be updated"
	}
}

func ParseAndExecute(messageArray []string) (string, bool) {
	Command := messageArray[0]
	var args []string
	if len(messageArray) > 1 {
		args = messageArray[1:]
	} else {
		args = []string{}
	}
	if handler, ok := commandParser[Command]; ok {
		output := handler(args)
		return output, ok
	}
	return "", false
}
