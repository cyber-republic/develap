package node

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List node nodes",
	Long:  `List node nodes`,
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("node list called with environment: [%s]\n\n", Env)

		containers := GetRunningContainersList()

		for _, container := range containers {
			for _, containerName := range container.Names {
				if strings.Contains(containerName, ContainerPrefix) && strings.Contains(containerName, Env) {
					i, err := strconv.ParseInt(strconv.FormatInt(container.Created, 10), 10, 64)
					if err != nil {
						log.Fatal(err)
					}
					created := time.Unix(i, 0)
					var portString string
					for _, port := range container.Ports {
						if port.IP == "0.0.0.0" {
							portString = fmt.Sprintf("%v", port.PublicPort)
						}
					}
					fmt.Printf("Name: %v\nID: %v\nImage: %v\nCmd: %v\nCreated: %v\nStatus: %v\nPort: %v\n\n",
						containerName[1:], container.ID[:10], container.Image, container.Command,
						created, container.Status, portString)
					break
				}
			}
		}
	},
}

func init() {
}
