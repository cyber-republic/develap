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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// RunCmd represents the run command
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run blockchain nodes",
	Long:  `Run blockchain nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("blockchain run called with environment: [%s] and nodes: [%s]\n", Env, Nodes)

		if !strings.EqualFold(Env, "mainnet") &&
			!strings.EqualFold(Env, "testnet")  {
			log.Fatalf("%s not recognized as a valid net type\n", Env)
		}

		ctx := context.Background()
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal(err)
		}

		nodes := strings.Split(strings.Replace(Nodes, " ", "", -1), ",")
			for _, node := range nodes {
				if node == "mainchain" || node == "did" || node == "eth" {
					if containerName, resp, err := runDockerContainer(ctx, cli, node); err != nil {
						log.Print(err)
					} else {
						fmt.Printf("\nNetwork: %s\nNode Type: %s\nContainer Name: %s\nContainer ID: %v\n",
										Env, node, containerName, resp.ID[:10])
					}
				} else {
					log.Fatalf("%s not recognized as a valid node type\n", node)
				}
			}
	},
}

func runDockerContainer(ctx context.Context, cli *client.Client, node string) (string, container.ContainerCreateCreatedBody, error) {
	var (
		containerName string
		resp container.ContainerCreateCreatedBody
		err  error
	)
	imageName := NodeDockerPath[node].ImageName
	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return containerName, resp, err
	}

	var containerRESTPort, containerRPCPort, hostRESTPort, hostRPCPort nat.Port
	if node == "mainchain" {
		containerRESTPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerRESTPort))
		containerRPCPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerRPCPort))
		hostRESTPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostRESTPort))
		hostRPCPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostRPCPort))
	} else if node == "did" {
		containerRESTPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerRESTPort))
		containerRPCPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerRPCPort))
		hostRESTPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostRESTPort))
		hostRPCPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostRPCPort))
	} else if node == "eth" {
		containerRPCPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerRPCPort))
		hostRPCPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostRPCPort))
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return containerName, resp, err
	}
	volumeData := filepath.FromSlash(fmt.Sprintf("%s/.data/%s/%s", currentDir, Env, node))
	os.MkdirAll(volumeData, os.ModePerm)
	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: volumeData,
			Target: NodeDockerPath[node].DataPath,
		},
	}

	portBindings := nat.PortMap {
		containerRPCPort:  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostRPCPort.Port()}},
	}

	containerConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{
			containerRPCPort:  struct{}{},
		},
	}

	if node == "eth" {
		if Env == "testnet" {
			containerConfig.Entrypoint = strslice.StrSlice{"/bin/sh"}
			containerConfig.Cmd = strslice.StrSlice{
				"-c", "./geth --testnet --datadir elastos_eth --gcmode 'archive' --rpc --rpcaddr 0.0.0.0 --rpccorsdomain '*' --rpcvhosts '*' --rpcport 20636 --rpcapi 'eth,net,web3' --ws --wsaddr 0.0.0.0 --wsorigins '*' --wsport 20635 --wsapi 'eth,net,web3'",
			}
		}
	} else {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: filepath.FromSlash(fmt.Sprintf("%s/%s/%s/config.json", currentDir, Env, node)),
			Target: NodeDockerPath[node].ConfigPath,
		})
		containerConfig.ExposedPorts[containerRESTPort] = struct{}{}
		portBindings[containerRESTPort] = []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostRESTPort.Port()}}
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Mounts: mounts,
	}

	containerName = fmt.Sprintf("develap-%s-%s-node", Env, node)

	resp, err = cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		containerName,
	)
	if err != nil {
		return containerName, resp, err
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return containerName, resp, err
	}
	return containerName, resp, nil
}

func init() {
	RunCmd.Flags().StringVarP(&Nodes, "nodes", "n", "", "Nodes to use [mainchain,did,eth]")
	RunCmd.MarkFlagRequired("nodes")
}
