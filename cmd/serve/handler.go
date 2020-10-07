package serve

import (
	"encoding/json"
	"fmt"
	"github.com/cyber-republic/develap/cmd/node"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	htmlIndex = `<html><body>Welcome!</body></html>`
)

type BasicContainerInfo struct {
	Environment string `json:"environment"`
	NodeType string `json:"node_type"`
	DockerImage string `json:"docker_image"`
	Endpoint string `json:"endpoint"`
	Status string `json:"status"`
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlIndex)
}

func HandleNodeEndpoints(mux *http.ServeMux) {
	containers := node.GetRunningContainersList()
	for _, container := range containers {
		for _, containerName := range container.Names {
			if strings.Contains(containerName, node.ContainerPrefix) {
				for _, port := range container.Ports {
					if port.IP == "0.0.0.0" {
						portString := fmt.Sprintf("%v", port.PublicPort)
						urlToParse := fmt.Sprintf("http://localhost:%s", portString)
						remoteURL, err := url.Parse(urlToParse)
						if err != nil {
							panic(err)
						}

						nodeType := strings.Split(containerName, "-")[2]
						localURL := fmt.Sprintf("/%s/%s", strings.Split(containerName, "-")[1], nodeType)

						proxy := httputil.NewSingleHostReverseProxy(remoteURL)
						mux.HandleFunc(localURL, rProxyHandler(proxy))
					}
				}
				break
			}
		}
	}
}

func HandleStatusAllNodesEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	runningNodes := node.GetRunningContainerInfo()

	var stoppedNodes []BasicContainerInfo

	for nodeType, dockerPath := range node.NodeDockerPath {
		for env, _ := range dockerPath.PortMapping {
			containerName := fmt.Sprintf("%s-%s-%s-node", node.ContainerPrefix, env, nodeType)
			found := false
			for _, container := range runningNodes {
				if container.ContainerName == containerName {
					found = true
					node := BasicContainerInfo{
						env,
						nodeType,
						dockerPath.ImageName,
						fmt.Sprintf("/%s/%s", strings.Split(containerName, "-")[1], nodeType),
						"Running",
					}
					stoppedNodes = append(stoppedNodes, node)
					break
				}
			}
			if !found {
				node := BasicContainerInfo{
					env,
					nodeType,
					dockerPath.ImageName,
					fmt.Sprintf("/%s/%s", strings.Split(containerName, "-")[1], nodeType),
					"Not running",
				}
				stoppedNodes = append(stoppedNodes, node)
			}
		}
	}
	json.NewEncoder(w).Encode(stoppedNodes)
}

func HandleStatusRunningNodesEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	runningNodes := node.GetRunningContainerInfo()
	json.NewEncoder(w).Encode(runningNodes)
}

func HandleStatusStoppedNodesEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	runningNodes := node.GetRunningContainerInfo()

	var stoppedNodes []BasicContainerInfo

	for nodeType, dockerPath := range node.NodeDockerPath {
		for env, _ := range dockerPath.PortMapping {
			containerName := fmt.Sprintf("%s-%s-%s-node", node.ContainerPrefix, env, nodeType)
			found := false
			for _, container := range runningNodes {
				if container.ContainerName == containerName {
					found = true
					break
				}
			}
			if !found {
				node := BasicContainerInfo{
					env,
					nodeType,
					dockerPath.ImageName,
					fmt.Sprintf("/%s/%s", strings.Split(containerName, "-")[1], nodeType),
					"Not running",
				}
				stoppedNodes = append(stoppedNodes, node)
			}
		}
	}
	json.NewEncoder(w).Encode(stoppedNodes)
}