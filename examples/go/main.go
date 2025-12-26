package main

import (
	sdk "emberdb/sdk/go"
	"fmt"
)

func main() {
	client := sdk.Connect("http://localhost:9182")
	resp, err := client.SetKey("lmfao", "an internet slang")
	if err != nil {
		fmt.Println(err)
	}
	resp, err = client.GetKey("lmfao")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Value)
}
