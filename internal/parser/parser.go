package parser

import "emberdb/internal/db"

var (
	commandParser = map[string]CommandFunc{
		"SET":    handleSet,
		"GET":    handleGet,
		"DEL":    handleDelete,
		"UPDATE": handleUpdate,
	}
)

type CommandFunc func(args []string) string

func handleSet(args []string) string {
	key := args[0]
	value := args[1]
	db.SetValue(key, value)
	return "SET OK"
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
	args := messageArray[1:]
	if handler, ok := commandParser[Command]; ok {
		output := handler(args)
		return output, ok
	}
	return "", false
}
