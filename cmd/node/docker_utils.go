package node

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
	"log"
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