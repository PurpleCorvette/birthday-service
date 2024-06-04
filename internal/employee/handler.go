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
		h.CreateEmployee(w, r)
	case "GET":
		h.GetEmployee(w, r)
	case "PUT":
		h.UpdateEmployee(w, r)
	case "DELETE":
		h.DeleteEmployee(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newEmp, err := h.service.AddEmployee(emp.Name, emp.DOB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEmp)
}

func (h *EmployeeHandler) GetEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	emp, err := h.service.GetEmployee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedEmp, err := h.service.UpdateEmployee(id, emp.Name, emp.DOB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedEmp)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid employee ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
