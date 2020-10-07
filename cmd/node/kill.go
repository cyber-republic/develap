package node

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
	Short: "Kill node nodes",
	Long:  `Kill node nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("node kill called with environment: [%s] and nodes: [%s]\n", Env, Type)

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal(err)
		}

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			log.Fatal(err)
		}

		nodes := strings.Split(strings.Replace(Type, " ", "", -1), ",")
		for _, container := range containers {
			for _, containerName := range container.Names {
				if strings.Contains(containerName, ContainerPrefix) && strings.Contains(containerName, Env) {
					if len(nodes) == 0 {
						log.Printf("Stopping container '%v' with ID '%v'...\n", containerName[1:], container.ID[:10])
						if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
							log.Fatal(err)
						}
						if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
							log.Fatal(err)
						}
					} else {
						for _, node := range nodes {
							if strings.Contains(containerName, node) {
								log.Printf("Stopping container '%v' with ID '%v'...\n", containerName[1:], container.ID[:10])
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
	},
}

func init() {
	usage := fmt.Sprintf("Type of node %v", SupportedNodes)
	KillCmd.Flags().StringVarP(&Type, "type", "t", "", usage)
}
