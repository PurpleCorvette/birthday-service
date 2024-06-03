package employee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	service EmployeeService
}

func NewEmployeeHandler(service EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var emp Employee
		if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		newEmp, err := h.service.AddEmployee(emp.Name, emp.DOB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newEmp)

	case "GET":
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "invalid employee ID", http.StatusBadRequest)
			return
		}

		emp, err := h.service.GetEmployee(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(emp)

	case "PUT":
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "invalid employee ID", http.StatusBadRequest)
			return
		}

		var emp Employee
		if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
			http.Error(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		updateEmp, err := h.service.UpdateEmployee(id, emp.Name, emp.DOB)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updateEmp)
	case "DELETE":
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			http.Error(w, "invalid employee id", http.StatusBadRequest)
			return
		}

		if err := h.service.DeleteEmployee(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
