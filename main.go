package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func newApiConfig() *apiConfig {
	return &apiConfig{fileserverHits: 0}
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	apiCfg := newApiConfig()
	mux.Handle("GET /app/", apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app/", http.FileServer(http.Dir("."))),
	))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.getMetrics)
	mux.HandleFunc("GET /api/reset", apiCfg.resetMetrics)

	corsMux := middlewareCors(mux)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
