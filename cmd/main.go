package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"birthday-service/internal/auth"
	"birthday-service/internal/employee"
	"birthday-service/pkg/logging"
)

func main() {
	log := logging.GetLogger()

	r := mux.NewRouter()

	authService := auth.NewAuthService()
	authHandler := auth.NewAuthHandler(authService)

	employeeService := employee.NewEmployeeService()
	employeeHandler := employee.NewEmployeeHandler(employeeService)

	r.Handle("/auth", authHandler).Methods("POST", "GET")
	r.Handle("/employee/{id:[0-9]+}", employeeHandler).Methods("GET", "PUT", "DELETE")
	r.Handle("/employee", employeeHandler).Methods("POST")

	log.Infoln("Starting server on :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
