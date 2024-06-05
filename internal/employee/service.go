package employee

import (
	"context"
	"time"

	"birthday-service/internal/db"
)

type Employee struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Birthday time.Time `json:"birthday"`
}

type EmployeeService interface {
	AddEmployee(name string, birthday time.Time) (Employee, error)
	UpdateEmployee(id int, name string, birthday time.Time) (Employee, error)
	DeleteEmployee(id int) error
	GetEmployee(id int) (Employee, error)
	GetAllEmployees() ([]Employee, error)
}

type employeeService struct {
	db db.DB
}

func NewEmployeeService(db db.DB) EmployeeService {
	return &employeeService{db: db}
}

func (s *employeeService) AddEmployee(name string, birthday time.Time) (Employee, error) {
	var employee Employee
	err := s.db.QueryRow(context.Background(),
		"INSERT INTO employees (name, birthday) VALUES ($1, $2) RETURNING id, name, birthday",
		name, birthday).Scan(&employee.ID, &employee.Name, &employee.Birthday)
	if err != nil {
		return Employee{}, err
	}

	return employee, nil
}

func (s *employeeService) UpdateEmployee(id int, name string, birthday time.Time) (Employee, error) {
	var employee Employee
	err := s.db.QueryRow(context.Background(),
		"UPDATE employees SET name=$1, birthday=$2 WHERE id=$3 RETURNING id, name, birthday",
		name, birthday, id).Scan(&employee.ID, &employee.Name, &employee.Birthday)
	if err != nil {
		return Employee{}, err
	}

	return employee, nil
}

func (s *employeeService) DeleteEmployee(id int) error {
	_, err := s.db.Exec(context.Background(),
		"DELETE FROM employees WHERE id=$1", id)
	return err
}

func (s *employeeService) GetEmployee(id int) (Employee, error) {
	var employee Employee
	err := s.db.QueryRow(context.Background(),
		"SELECT id, name, birthday FROM employees WHERE id=$1",
		id).Scan(&employee.ID, &employee.Name, &employee.Birthday)
	if err != nil {
		return Employee{}, err
	}

	return employee, nil
}

func (s *employeeService) GetAllEmployees() ([]Employee, error) {
	rows, err := s.db.Query(context.Background(),
		"SELECT id, name, birthday FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Birthday)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}
