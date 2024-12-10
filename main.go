package main

import (
	"net/http"
	"sync/atomic"
)

func appIndexHandler() http.Handler {
	return http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
}

func redirectToApp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
}

func main() {
	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
	}
	mux := http.NewServeMux()

	//App Handlers
	mux.HandleFunc("/", redirectToApp)
	mux.Handle("/app/", apiCfg.middlewareMetricInc(appIndexHandler()))

	//Admin Handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	//API Handlers
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
