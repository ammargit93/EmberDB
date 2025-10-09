package main

import (
	"fmt"

	"emberdb/client/sdk"
)

func main() {
	client := sdk.Connect("localhost:1010")
	val, _ := client.GetValue("a")
	fmt.Println(val)

}
