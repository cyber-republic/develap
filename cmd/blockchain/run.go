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

	var containerPort, hostPort nat.Port
	if node == "mainchain" {
		containerPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerPort))
		hostPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostPort))
	} else if node == "did" {
		containerPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerPort))
		hostPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostPort))
	} else if node == "eth" {
		containerPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].ContainerPort))
		hostPort = nat.Port(fmt.Sprintf("%s/tcp", NodeDockerPath[node].PortMapping[Env].HostPort))
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
		containerPort: []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: hostPort.Port()}},
	}

	containerConfig := &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
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
			Source: filepath.FromSlash(fmt.Sprintf("%s/node_config/%s/%s.config.json", currentDir, Env, node)),
			Target: NodeDockerPath[node].ConfigPath,
		})
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
