package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	Short: "Setup route to different docker containers",
	Long:  `Setup route to different docker containers`,
	Run: func(cmd *cobra.Command, args []string) {
		var httpSrv *http.Server

		remote, err := url.Parse("http://localhost:21336")
		if err != nil {
			panic(err)
		}

		mux := &http.ServeMux{}
		mux.HandleFunc("/", handleIndex)
		proxy := httputil.NewSingleHostReverseProxy(remote)
		mux.HandleFunc("/testnet/blockchain/mainchain", rProxyHandler(proxy))
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

		err = httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
