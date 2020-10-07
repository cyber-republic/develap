package cmd

import (
	"github.com/cyber-republic/develap/cmd/serve"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Setup route to different node containers",
	Long:  `Setup route to different node containers`,
	Run: func(c *cobra.Command, args []string) {
		var httpSrv *http.Server

		mux := &http.ServeMux{}

		// Default route
		mux.HandleFunc("/", serve.HandleIndex)

		// Node routes
		serve.HandleNodeEndpoints(mux)

		// Status routes
		mux.HandleFunc("/status/nodes", serve.HandleStatusAllNodesEndpoint)
		mux.HandleFunc("/status/nodes/running", serve.HandleStatusRunningNodesEndpoint)
		mux.HandleFunc("/status/nodes/stopped", serve.HandleStatusStoppedNodesEndpoint)

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
		log.Println("Starting server http://localhost:5000")

		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
