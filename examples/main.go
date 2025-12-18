package main

import (
	sdk "emberdb/sdk/go"
	"fmt"
)

func main() {
	client := sdk.Connect("http://localhost:9182")
	resp, err := client.SetKey("age", 101)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Key, resp.Value)
}
