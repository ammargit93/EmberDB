package main

import (
	sdk "emberdb/sdk/go"
	"fmt"
)

func main() {
	client := sdk.Connect("http://localhost:9182")

	values, err := client.MGetKey([]string{"a", "b"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(values["a"]) // 10
	fmt.Println(values["b"]) // 20
}
