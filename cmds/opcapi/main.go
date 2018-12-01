package main

// To start opcapi in powershell
// & {$ENV:OPC_SERVER="Graybox.Simulator"; $ENV:OPC_NODES="localhost";  go run main.go -addr ":8765"}

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/konimarti/opc"
	"github.com/konimarti/opc/api"
)

var (
	addr    = flag.String("addr", ":8765", "enter address to start api")
	cfgFile = flag.String("config", "opcapi.conf", "config file name")
)

type tmlConfig struct {
	Config api.Config `toml:"config"`
	Opc    OpcConfig  `toml:"opc"`
}

type OpcConfig struct {
	Server string `toml:"server"`
	Nodes  []string
	Tags   []string
}

func main() {
	flag.Parse()

	opc.Debug()

	// parse config
	data, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	// parse config
	var cfg tmlConfig
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		log.Fatal(err)
	}

	server := cfg.Opc.Server
	if server == "" {
		server = strings.Trim(os.Getenv("OPC_SERVER"), " ")
		if server == "" {
			panic("OPC_SERVER not set")
		}
	}
	nodes := cfg.Opc.Nodes
	if len(nodes) == 0 {
		nodes = strings.Split(os.Getenv("OPC_NODES"), ",")
		if len(nodes) == 0 {
			panic("OPC_NODES not set; separate nodes with ','")
		}
	}
	for i := range nodes {
		nodes[i] = strings.Trim(nodes[i], " ")
	}

	fmt.Println("API starting with OPC", server, nodes, *addr)

	client := opc.NewConnection(
		server,
		nodes,
		cfg.Opc.Tags,
	)
	defer client.Close()

	app := api.App{Config: cfg.Config}
	app.Initialize(client)

	app.Run(*addr)
}
