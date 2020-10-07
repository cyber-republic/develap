package serve

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func rProxyHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL)
		w.Header().Set("X-GoProxy", "GoProxy")
		p.ServeHTTP(w, r)
	}
}
