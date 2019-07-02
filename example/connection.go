package main

import (
	"fmt"

	"github.com/konimarti/opc"
)

func main() {
	client, err := opc.NewConnection(
		"Graybox.Simulator",
		[]string{"localhost"},
		[]string{"numeric.sin.int64", "numeric.saw.float"},
	)
	defer client.Close()

	if err != nil {
		panic(err)
	}

	// read single tag: value, quality, timestamp
	fmt.Println(client.ReadItem("numeric.sin.int64"))

	// read all added tags
	fmt.Println(client.Read())
}
