package main

import (
	"fmt"

	"github.com/konimarti/opc"
	"github.com/spf13/cobra"
)

func main() {
	//var server, node string

	var cmdList = &cobra.Command{
		Use:   "list [node]",
		Short: "Lists the OPC servers on node.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			node := args[0]
			servers_found := opc.NewAutomationObject().GetOPCServers(node)
			fmt.Printf("Found %d server(s) on '%s':\n", len(servers_found), node)
			for _, server := range servers_found {
				fmt.Println(server)
			}
		},
	}

	var rootCmd = &cobra.Command{Use: "opc-cli"}
	//rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "Graybox.Simulator", "ProgID for OPC Server")
	//rootCmd.PersistentFlags().StringVarP(&node, "node", "n", "localhost", "Node for OPC Server")

	rootCmd.AddCommand(cmdList)
	rootCmd.Execute()
}
