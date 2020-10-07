package cmd

import (
	"fmt"
	"github.com/cyber-republic/develap/cmd/node"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	htmlIndex = `<html><body>Welcome!</body></html>`
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlIndex)
}

func rProxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-GoProxy", "GoProxy")
		p.ServeHTTP(w, r)
	}
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Setup route to different node containers",
	Long:  `Setup route to different node containers`,
	Run: func(c *cobra.Command, args []string) {
		var httpSrv *http.Server

		mux := &http.ServeMux{}
		mux.HandleFunc("/", handleIndex)

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

							nodeName := strings.Split(containerName, "-")[2]
							var localURL string
							if strings.Contains(containerName, node.TestNet) {
								localURL = fmt.Sprintf("/%s/%s", node.TestNet, nodeName)
							} else if strings.Contains(containerName, node.MainNet) {
								localURL = fmt.Sprintf("/%s/%s", node.MainNet, nodeName)
							}

							proxy := httputil.NewSingleHostReverseProxy(remoteURL)
							fmt.Printf("Remote: %s Local: %s\n", remoteURL, localURL)
							mux.HandleFunc(localURL, rProxyHandler(proxy))
						}
					}
					break
				}
			}
		}

		// set timeouts so that a slow or malicious client doesn't
		// hold resources forever
		httpSrv = &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      mux,
		}
		httpSrv.Addr = ":5000"

		// Launch HTTP server
		fmt.Println("Starting server http://localhost:5000")

		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	},
}

func blockchainEndpoints() {

}

func init() {
	rootCmd.AddCommand(serveCmd)
}
