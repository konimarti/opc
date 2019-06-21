package main

import (
	"fmt"
	"os"

	"github.com/konimarti/opc"
	"github.com/spf13/cobra"
)

func main() {

	var cmdList = &cobra.Command{
		Use:   "list [node]",
		Short: "Lists the OPC servers available on a specific node.",
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
		Use:   "info [node] [server]",
		Short: "Try connect the OPC server on a specific node and check if it is running.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			nodes := []string{args[0]}
			server := args[1]
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

	var cmdBrowse = &cobra.Command{
		Use:   "browse [node] [server] [branch_name]",
		Short: "Browse OPC tags. If only sub-branch is requested, use optional branch_name.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			nodes := []string{args[0]}
			server := args[1]
			name := "root"
			if len(args) > 2 {
				name = args[2]
			}
			tree, err := opc.CreateBrowser(server, nodes)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			opc.PrettyPrint(opc.ExtractBranchByName(tree, name))
		},
	}

	var cmdRead = &cobra.Command{
		Use:   "read [node] [server] [tags...]",
		Short: "Read OPC tags.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			nodes := []string{args[0]}
			server := args[1]
			tags := args[2:]
			conn := opc.NewConnection(
				server,
				nodes,
				tags,
			)
			fmt.Println(conn.Read())
		},
	}

	var cmdWrite = &cobra.Command{
		Use:   "write [node] [server] [tag] [value]",
		Short: "Write value to OPC tag.",
		Long:  ``,
		Args:  cobra.MinimumNArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			nodes := []string{args[0]}
			server := args[1]
			tag := args[2]
			value := args[3]
			conn := opc.NewConnection(
				server,
				nodes,
				[]string{tag},
			)
			conn.Write(tag, value)
		},
	}

	var rootCmd = &cobra.Command{Use: "opc-cli"}

	rootCmd.AddCommand(cmdList, cmdInfo, cmdBrowse, cmdRead, cmdWrite)
	rootCmd.Execute()
}
