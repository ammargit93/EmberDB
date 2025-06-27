package api

import (
	"fmt"
	"strconv"
	"strings"
)

func SplitIPandIncrementPort(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) < 3 {
		return "http://localhost:6060"
	}
	host := parts[1]    // "//localhost"
	portStr := parts[2] // "9090"
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		return "http://localhost:6060"
	}
	newPort := portInt + 1
	host = strings.TrimPrefix(host, "//")
	return fmt.Sprintf("http://%s:%d", host, newPort)
}

func IncrementPort(portStr string) string {
	portNum := strings.TrimPrefix(portStr, ":")
	portInt, err := strconv.Atoi(portNum)
	if err != nil {
		return ":6060"
	}

	return fmt.Sprintf(":%d", portInt+1)
}
