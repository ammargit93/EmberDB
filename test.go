package main

import (
	"emberdb/sdk"
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	sdk.Connect("localhost:1010")

	err := sdk.SetValue("a", "10")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
	fmt.Println("Hello World")
}
