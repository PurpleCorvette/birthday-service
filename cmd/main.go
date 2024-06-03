package main

import (
	"net/http"

	"birthday-service/internal/auth"
	"birthday-service/pkg/logging"

	"github.com/gorilla/mux"
)

func main() {
	log := logging.GetLogger()

	r := mux.NewRouter()

	authService := auth.NewAuthService()
	authHandler := auth.NewAuthHandler(authService)

	r.Handle("/auth", authHandler).Methods("POST", "GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
