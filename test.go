package main

import (
	"fmt"

	"emberdb/client/sdk"
)

func main() {
	client := sdk.Connect("localhost:1010")

	_ = client.SetFile("a", "C:\\Users\\Ammar1\\go\\emberdb\\README.md")

	val, _ := client.GetFile("a")
	fmt.Println(val)

}
