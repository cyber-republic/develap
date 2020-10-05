package cmd

import (
	"os"

	"github.com/cyber-republic/develap/cmd/blockchain"
	"github.com/spf13/cobra"
)

// blockchainCmd represents the blockchain command
var blockchainCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Interact with blockchain nodes",
	Long:  `Interact with blockchain nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func init() {
	blockchainCmd.AddCommand(blockchain.ListCmd)
	blockchainCmd.AddCommand(blockchain.RunCmd)
	blockchainCmd.AddCommand(blockchain.KillCmd)
	rootCmd.AddCommand(blockchainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blockchainCmd.PersistentFlags().String("foo", "", "A help for foo")
	blockchainCmd.PersistentFlags().StringVarP(&blockchain.Env, "env", "e", "", "environment to use [mainnet,testnet]")
	blockchainCmd.MarkFlagRequired("env")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//blockchainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
