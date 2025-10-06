// internal/parser/parser.go
package parser

import (
	"fmt"
	"strings"

	"emberdb/internal/db"
)

var (
	commandParser = map[string]CommandFunc{
		"SET":      handleSet,
		"GET":      handleGet,
		"DEL":      handleDelete,
		"UPDATE":   handleUpdate,
		"SETFILE":  setFile,
		"GETFILE":  getFile,
		"SAVEFILE": saveFile,
	}
)

type CommandFunc func(args []string) string

func getFile(args []string) string {
	if len(args) < 1 {
		return "Missing key for GETFILE"
	}
	b, err := db.GetFile(args[0])
	if err != nil {
		return "Error getting file: " + err.Error()
	}
	return b
}

func setFile(args []string) string {
	if len(args) < 2 {
		return "Missing arguments for SETFILE"
	}
	if err, ok := db.SetFile(args[0], args[1]); err != nil || !ok {
		return "Error setting the file"
	}
	return "File successfully set"
}

func saveFile(args []string) string {
	if len(args) < 2 {
		return "Missing arguments for SAVEFILE"
	}
	if err := db.SaveFile(args[0], args[1]); err != nil {
		return "Error Saving the file"
	}
	return "File Successfully saved to " + args[1]

}

func handleSet(args []string) string {
	if len(args) < 2 {
		return "Missing key or value for SET"
	}
	key := args[0]
	value := strings.Join(args[1:], " ")
	if db.SetValue(key, value) {
		return "Key already exists"
	}
	return "SET OK"
}

func handleGet(args []string) string {
	if len(args) < 1 {
		return "Missing key for GET"
	}
	val := db.GetValue(args[0])
	if str, ok := val.(string); ok {
		return str
	} else if b, ok := val.([]byte); ok {
		return string(b)
	}
	return fmt.Sprint(val)
}

func handleDelete(args []string) string {
	if len(args) < 1 {
		return "Missing key for DEL"
	}
	if db.DeleteKey(args[0]) {
		return "Key Deleted Successfully"
	}
	return "Key could not be deleted"
}

func handleUpdate(args []string) string {
	if len(args) < 2 {
		return "Missing key or value for UPDATE"
	}
	if db.UpdateValue(args[0], args[1]) {
		return "Key Updated Successfully"
	}
	return "Key could not be updated"
}

func ParseAndExecute(messageArray []string) (string, bool) {
	if len(messageArray) == 0 {
		return "Invalid command", false
	}
	command := strings.ToUpper(messageArray[0])
	args := messageArray[1:]
	if handler, ok := commandParser[command]; ok {
		return handler(args), true
	}
	return "Unknown command", false
}
