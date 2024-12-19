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
	mux.HandleFunc("PUT /api/users", apiCfg.middlewareAuth(apiCfg.handlerUpdateUser))

	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetOneChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.middlewareAuth(apiCfg.handlerCreateChirp))
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.middlewareAuth(apiCfg.handlerDeleteChirp))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Printf("Running chirpy server on http://localhost%s", server.Addr)
	server.ListenAndServe()
}
