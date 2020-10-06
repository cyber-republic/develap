package blockchain

import (
	"fmt"
	"net"
	"os"

	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
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

type DockerContainerDataDir struct {
	HostCreate    bool
	ContainerPath string
}

type DockerContainer struct {
	ContainerName         string
	ImageName             string
	Volumes               map[string]DockerContainerDataDir
	ContainerExposedPorts nat.PortSet
	HostPortMappings      nat.PortMap
	EntryPoint            strslice.StrSlice
	Cmd                   strslice.StrSlice
}

var (
	Env            string
	Nodes          string
	CurrentDir     = getCurrentDir()
	NetworkName    = "elastos_sws"
	NodeDockerPath = map[string]DockerPath{
		"mainchain": {
			ImageName:    "cyberrepublic/ela-mainchain:v0.6.0",
			DataPath:     "/elastos/elastos",
			ConfigPath:   "/elastos/config.json",
			PortMapping: map[string]DockerPort{
				"mainnet": {ContainerPort: "20336", HostPort: getRandomPort()},
				"testnet": {ContainerPort: "21336", HostPort: getRandomPort()},
			},
		},
		"did": {
			ImageName:  "cyberrepublic/ela-did-sidechain:v0.2.0",
			DataPath:   "/elastos/elastos_did",
			ConfigPath: "/elastos/config.json",
			PortMapping: map[string]DockerPort{
				"mainnet":  {ContainerPort: "20606", HostPort: getRandomPort()},
				"testnet":  {ContainerPort: "21606", HostPort: getRandomPort()},
			},
		},
		"eth": {
			ImageName: "cyberrepublic/ela-eth-sidechain:v0.1.0",
			DataPath:  "/elastos/elastos_eth",
			PortMapping: map[string]DockerPort{
				"mainnet": {ContainerPort: "20636", HostPort: getRandomPort()},
				"testnet": {ContainerPort: "20636", HostPort: getRandomPort()},
			},
		},
	}
)

func getRandomPort() string {
	listener, err := net.Listen("tcp", ":0")
	defer listener.Close()
	if err != nil {
		panic(err)
	}
	portInt := listener.Addr().(*net.TCPAddr).Port
	portString := fmt.Sprintf("%d", portInt)
	return portString
}

func getCurrentDir() string {
	var currentDir string
	if pwd, err := os.Getwd(); err == nil {
		currentDir = pwd
	}
	return currentDir
}
