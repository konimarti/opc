// To start opcapi in powershell
// & {$ENV:OPC_SERVER="Graybox.Simulator"; $ENV:OPC_NODES="localhost";  go run main.go -addr ":8765"}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/konimarti/opc"
	"github.com/konimarti/opc/api"
)

var (
	addr = flag.String("addr", ":4000", "enter address to start api")
)

func main() {
	flag.Parse()

	opc.Debug()

	server := strings.Trim(os.Getenv("OPC_SERVER"), " ")
	if server == "" {
		panic("OPC_SERVER not set")
	}
	nodes := strings.Split(os.Getenv("OPC_NODES"), ",")
	if len(nodes) == 0 {
		panic("OPC_NODES not set; separate nodes with ','")
	}
	for i, _ := range nodes {
		nodes[i] = strings.Trim(nodes[i], " ")
	}

	fmt.Println("API starting with OPC", server, nodes, *addr)

	client := opc.NewConnection(
		server,
		nodes,
		[]string{},
	)
	defer client.Close()

	app := api.App{}
	app.Initialize(client)

	app.Run(*addr)
}
