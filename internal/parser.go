package internal

import (
	"strings"
	"time"
)

var DurationMap map[string]time.Duration = map[string]time.Duration{
	"s": time.Second,
	"m": time.Second * 60,
	"h": time.Second * 60 * 60,
}

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
