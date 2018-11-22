package main

import (
	"fmt"
	"os"

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

	var cmdInfo = &cobra.Command{
		Use:   "info [server] [node]",
		Short: "Try connect the OPC servers on node and get some info.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			server := args[0]
			nodes := []string{args[1]}
			obj := opc.NewAutomationObject()
			_, err := obj.TryConnect(server, nodes)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if obj.IsConnected() {
				fmt.Printf("%s on '%v' is up and running.\n", server, nodes[0])
			} else {
				fmt.Printf("%s on '%v' is not running.\n", server, nodes[0])
			}

		},
	}

	var rootCmd = &cobra.Command{Use: "opc-cli"}
	//rootCmd.PersistentFlags().StringVarP(&server, "server", "s", "Graybox.Simulator", "ProgID for OPC Server")
	//rootCmd.PersistentFlags().StringVarP(&node, "node", "n", "localhost", "Node for OPC Server")

	rootCmd.AddCommand(cmdList, cmdInfo)
	rootCmd.Execute()
}
