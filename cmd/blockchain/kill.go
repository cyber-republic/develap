package blockchain

import (
	"fmt"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// KillCmd represents the kill command
var KillCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill blockchain nodes",
	Long:  `Kill blockchain nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("blockchain kill called with environment: [%s] and nodes: [%s]\n", Env, Nodes)

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal(err)
		}

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			log.Fatal(err)
		}

		nodes := strings.Split(strings.Replace(Nodes, " ", "", -1), ",")
		for _, container := range containers {
			for _, containerName := range container.Names {
				if strings.Contains(containerName, "develap") && strings.Contains(containerName, Env) {
					if len(nodes) == 0 {
						fmt.Printf("Stopping container '%v' with ID '%v'...\n", containerName[1:], container.ID[:10])
						if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
							log.Fatal(err)
						}
						if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
							log.Fatal(err)
						}
					} else {
						for _, node := range nodes {
							if strings.Contains(containerName, node) {
								fmt.Printf("Stopping container '%v' with ID '%v'...\n", containerName[1:], container.ID[:10])
								if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
									log.Fatal(err)
								}
								if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
									log.Fatal(err)
								}
							}
						}
					}
					break
				}
			}
		}

		networks, err := cli.NetworkList(ctx, types.NetworkListOptions{})
		if err != nil {
			log.Fatal(err)
		}
		for _, network := range networks {
			if network.Name == NetworkName {
				fmt.Printf("\nRemoving network '%v' with ID '%v'...\n", network.Name, network.ID)
				_ = cli.NetworkRemove(ctx, network.ID)
			}
		}
	},
}

func init() {
	KillCmd.Flags().StringVarP(&Nodes, "nodes", "n", "", "Nodes to use [mainchain,did,eth]")
}
