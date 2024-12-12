package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func appIndexHandler() http.Handler {
	return http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
}

func redirectToApp(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
}

func main() {
	apiCfg := setupConfig()
	mux := http.NewServeMux()

	//App Handlers
	mux.HandleFunc("/", redirectToApp)
	mux.Handle("/app/", apiCfg.middlewareMetricInc(appIndexHandler()))

	//Admin Handlers
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetricHits)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	//API Handlers
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Running chirpy server on http://localhost%s", server.Addr)
	server.ListenAndServe()
}
