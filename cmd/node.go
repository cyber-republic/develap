package cmd

import (
	"fmt"
	"os"

	"github.com/cyber-republic/develap/cmd/node"
	"github.com/spf13/cobra"
)

// nodeCmd represents the node command
var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Interact with nodes",
	Long:  `Interact with nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	nodeCmd.AddCommand(node.ListCmd)
	nodeCmd.AddCommand(node.RunCmd)
	nodeCmd.AddCommand(node.KillCmd)
	rootCmd.AddCommand(nodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")
	usage := fmt.Sprintf("Type of environment [%s %s]", node.MainNet, node.TestNet)
	nodeCmd.PersistentFlags().StringVarP(&node.Env, "env", "e", "", usage)
	nodeCmd.MarkFlagRequired("env")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
