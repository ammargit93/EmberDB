package internal

import (
	"strings"
)

func buildMapFromArgs(args *[]string) map[string]string {
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
	return ArgMap

}

func Parse(args []string) {
	_ = buildMapFromArgs(&args)
}
