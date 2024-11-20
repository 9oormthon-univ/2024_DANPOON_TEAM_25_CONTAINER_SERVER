package main

import (
	"log"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("jenkins test2"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health_check", HealthCheck)
	srv := http.Server{
		Addr:    ":8082",
		Handler: mux,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
