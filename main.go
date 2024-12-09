package main

import (
	"net/http"
	"sync/atomic"
)

func appIndexHandler() http.Handler {
	return http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
}

func main() {
	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricInc(appIndexHandler()))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetricHits)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
