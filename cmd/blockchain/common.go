/*
Copyright © 2019 Cyber Republic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package blockchain

import (
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
				"mainnet": {ContainerPort: "20336", HostPort: "20336"},
				"testnet": {ContainerPort: "21336", HostPort: "21336"},
			},
		},
		"did": {
			ImageName:  "cyberrepublic/ela-did-sidechain:v0.2.0",
			DataPath:   "/elastos/elastos_did",
			ConfigPath: "/elastos/config.json",
			PortMapping: map[string]DockerPort{
				"mainnet":  {ContainerPort: "20606", HostPort: "20606"},
				"testnet":  {ContainerPort: "21606", HostPort: "21606"},
			},
		},
		"eth": {
			ImageName: "cyberrepublic/ela-eth-sidechain:v0.1.0",
			DataPath:  "/elastos/elastos_eth",
			PortMapping: map[string]DockerPort{
				"mainnet": {ContainerPort: "20636", HostPort: "20636"},
				"testnet": {ContainerPort: "20636", HostPort: "21636"},
			},
		},
	}
)

func getCurrentDir() string {
	var currentDir string
	if pwd, err := os.Getwd(); err == nil {
		currentDir = pwd
	}
	return currentDir
}
