package node

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"time"
)

type DockerPort struct {
	ContainerPort     string
	HostPort          string
}

type DockerPath struct {
	ImageName    string
	DataPath     string
	ConfigPath   string
	PortMapping  map[string]DockerPort
}

type DockerDataDir struct {
	HostCreate    bool
	ContainerPath string
}

type DockerContainer struct {
	ContainerName         string
	ImageName             string
	Volumes               map[string]DockerDataDir
	ContainerExposedPorts nat.PortSet
	HostPortMappings      nat.PortMap
	EntryPoint            strslice.StrSlice
	Cmd                   strslice.StrSlice
}

type ContainerInfo struct {
	Environment string `json:"environment"`
	NodeType string `json:"node_type"`
	ContainerName string `json:"container_name"`
	ContainerID string `json:"container_id"`
	ContainerCmd string `json:"container_cmd"`
	DockerImage string `json:"docker_image"`
	Endpoint string `json:"endpoint"`
	Created time.Time `json:"created"`
	Status string `json:"status"`
	Port string `json:"port"`
}

func GetRunningContainerInfo() []ContainerInfo {
	var status []ContainerInfo
	containers := GetRunningContainersList()
	for _, container := range containers {
		for _, containerName := range container.Names {
			if strings.Contains(containerName, ContainerPrefix) {
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
				environment := strings.Split(containerName[1:], "-")[1]
				nodeType := strings.Split(containerName[1:], "-")[2]
				endpoint := fmt.Sprintf("/%s/%s", environment, nodeType)
				s := ContainerInfo{
					environment,
					nodeType,
					containerName[1:],
					container.ID[:10],
					container.Image,
					container.Command,
					endpoint,
					created,
					container.Status,
					portString,
				}
				status = append(status, s)
				break
			}
		}
	}
	return status
}

func IsSupportedNode(node string) bool {
	for _, n := range SupportedNodes {
		if n == node {
			return true
		}
	}
	return false
}

func GetRunningContainersList() []types.Container {
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