package blockchain

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blockchain nodes",
	Long:  `List blockchain nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("blockchain list called with environment: [%s]\n\n", Env)

		containers := GetContainersList()

		for _, container := range containers {
			for _, containerName := range container.Names {
				if strings.Contains(containerName, "develap") && strings.Contains(containerName, Env) {
					i, err := strconv.ParseInt(strconv.FormatInt(container.Created, 10), 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					created := time.Unix(i, 0)
					ports := make(map[string]string, 0)
					for _, port := range container.Ports {
						if port.IP == "0.0.0.0" {
							portString := fmt.Sprintf("%v", port.PublicPort)
							if strings.Contains(containerName, "mainchain") {
								ports[portString] = "RPC"
							} else if strings.Contains(containerName, "did") {
								ports[portString] = "RPC"
							} else if strings.Contains(containerName, "eth") {
								ports[portString] = "RPC"
							}
						}
					}
					fmt.Printf("Name: %v\nID: %v\nImage: %v\nCmd: %v\nCreated: %v\nStatus: %v\nPorts: %v\n\n",
						containerName[1:], container.ID[:10], container.Image, container.Command,
						created, container.Status, ports)
					break
				}
			}
		}
	},
}

func GetContainersList() []types.Container {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return containers
}

func init() {
}
