package main

import (
	"fmt"

	"github.com/konimarti/opc"
)

func main() {
	progid := "Graybox.Simulator"
	nodes := []string{"localhost"}

	// create browser and collect all tags on OPC server
	browser, err := opc.CreateBrowser(progid, nodes)
	if err != nil {
		panic(err)
	}

	// extract subtree
	subtree := opc.ExtractBranchByName(browser, "textual")

	// print out all the information
	opc.PrettyPrint(subtree)

	// create opc connection with all tags from subtree
	conn := opc.NewConnection(
		progid,
		nodes,
		opc.CollectTags(subtree),
	)
	defer conn.Close()

	fmt.Println(conn.Read())
}
