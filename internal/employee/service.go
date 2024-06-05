package employee

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"birthday-service/internal/db"
)

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

type EmployeeService interface {
	AddEmployee(name, dob string) (Employee, error)
	GetEmployee(id int) (Employee, error)
	UpdateEmployee(id int, name, dob string) (Employee, error)
	DeleteEmployee(id int) error
}

type employeeService struct{}

func NewEmployeeService() EmployeeService {
	return &employeeService{}
}

func (s *employeeService) AddEmployee(name, dob string) (Employee, error) {
	var emp Employee
	err := db.Conn.QueryRow(context.Background(),
		"INSERT INTO employees (name, dob) VALUES ($1, $2) RETURNING id, name, dob",
		name, dob).Scan(&emp.ID, &emp.Name, &emp.DOB)
	if err != nil {
		return Employee{}, err
	}

	return emp, nil
}

func (s *employeeService) GetEmployee(id int) (Employee, error) {
	var emp Employee
	err := db.Conn.QueryRow(context.Background(),
		"SELECT id, name, dob FROM employees WHERE id=$1", id).Scan(&emp.ID, &emp.Name, &emp.DOB)
	if err == pgx.ErrNoRows {
		return Employee{}, errors.New("employee not found")
	}
	if err != nil {
		return Employee{}, err
	}

	return emp, nil
}

func (s *employeeService) UpdateEmployee(id int, name, dob string) (Employee, error) {
	var emp Employee
	err := db.Conn.QueryRow(context.Background(),
		"UPDATE employees SET name=$1, dob=$2 WHERE id=$3 RETURNING id, name, dob",
		name, dob, id).Scan(&emp.ID, &emp.Name, &emp.DOB)
	if err == pgx.ErrNoRows {
		return Employee{}, errors.New("employee not found")
	}
	if err != nil {
		return Employee{}, err
	}

	return emp, nil
}

func (s *employeeService) DeleteEmployee(id int) error {
	_, err := db.Conn.Exec(context.Background(), "DELETE FROM employees WHERE id=$1", id)
	if err == pgx.ErrNoRows {
		return errors.New("employee not found")
	}
	if err != nil {
		return err
	}

	return nil
}
