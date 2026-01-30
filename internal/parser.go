package internal

import (
	"strings"
	"sync"
)

var mu sync.RWMutex

func buildMapFromArgs(args *[]string) {
	mu.Lock()
	defer mu.Unlock()
	if ArgMap == nil {
		ArgMap = make(map[string]string)
	}
	for i, arg := range *args {
		var flag, value string
		if strings.HasPrefix(arg, "--") {
			flag = strings.Replace(arg, "--", "", 2)
			value = (*args)[i+1]
		}
		ArgMap[flag] = value
	}

}

func Parse(args []string) {
	buildMapFromArgs(&args)
}
