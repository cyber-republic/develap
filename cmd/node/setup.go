package node

import (
	"github.com/cyber-republic/develap/cmd/system"
)

const (
	ContainerPrefix = "sws"
	MainNet = "mainnet"
	TestNet = "testnet"
)

var (
	SupportedNodes = []string{
		"mainchain",
		"did",
		"eth",
	}

	Env            string
	Type           string
	CurrentDir     = system.GetCurrentDir()
	NodeDockerPath = map[string]DockerPath{
		"mainchain": {
			ImageName:    "cyberrepublic/ela-mainchain:v0.6.0",
			DataPath:     "/elastos/elastos",
			ConfigPath:   "/elastos/config.json",
			PortMapping: map[string]DockerPort{
				MainNet: {ContainerPort: "20336", HostPort: system.GetRandomPort()},
				TestNet: {ContainerPort: "21336", HostPort: system.GetRandomPort()},
			},
		},
		"did": {
			ImageName:  "cyberrepublic/ela-did-sidechain:v0.2.0",
			DataPath:   "/elastos/elastos_did",
			ConfigPath: "/elastos/config.json",
			PortMapping: map[string]DockerPort{
				MainNet: {ContainerPort: "20606", HostPort: system.GetRandomPort()},
				TestNet: {ContainerPort: "21606", HostPort: system.GetRandomPort()},
			},
		},
		"eth": {
			ImageName: "cyberrepublic/ela-eth-sidechain:v0.1.0",
			DataPath:  "/elastos/elastos_eth",
			PortMapping: map[string]DockerPort{
				MainNet: {ContainerPort: "20636", HostPort: system.GetRandomPort()},
				TestNet: {ContainerPort: "20636", HostPort: system.GetRandomPort()},
			},
		},
	}
)
