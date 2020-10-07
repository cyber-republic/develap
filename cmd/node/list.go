package node

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List node nodes",
	Long:  `List node nodes`,
	Run: func(c *cobra.Command, args []string) {
		log.Printf("node list called with environment: [%s]\n\n", Env)

		runningNodes := GetRunningContainerInfo()

		for _, container := range runningNodes {
			fmt.Printf("Environment: %v\nNode Type: %v\nContainer Name: %v\nContainer ID: %v\nContainer Cmd: %v\n" +
				"Docker Image: %v\nEndpoint: %v\nCreated: %v\nStatus: %v\nPort: %v\n\n",
				container.Environment, container.NodeType,
				container.ContainerName, container.ContainerID, container.ContainerCmd, container.DockerImage,
				container.Endpoint, container.Created, container.Status, container.Port)
		}
	},
}

func init() {
}
