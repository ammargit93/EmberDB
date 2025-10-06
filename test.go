package main

import (
	"fmt"

	"emberdb/client/sdk"
)

func main() {
	// fmt.Println("Hello World")
	client := sdk.Connect("localhost:1010")
	val, _ := client.GetValue("a")
	fmt.Println(val)

}
