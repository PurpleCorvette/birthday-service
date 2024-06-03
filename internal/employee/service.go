package employee

import "errors"

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"` // day of birth
}

type EmployeeService interface {
	AddEmployee(name, dob string) (Employee, error)
	GetEmployee(id int) (Employee, error)
	UpdateEmployee(id int, name, dob string) (Employee, error)
	DeleteEmployee(id int) error
	ListEmployees() ([]Employee, error)
}

type employeeService struct {
	employees []Employee
	nextID    int
}

func NewEmployeeService() EmployeeService {
	return &employeeService{
		employees: []Employee{},
		nextID:    1,
	}
}

func (s *employeeService) AddEmployee(name, dob string) (Employee, error) {
	employee := Employee{
		ID:   s.nextID,
		Name: name,
		DOB:  dob,
	}
	s.nextID++
	s.employees = append(s.employees, employee)
	return employee, nil
}

func (s *employeeService) GetEmployee(id int) (Employee, error) {
	for _, emp := range s.employees {
		if emp.ID == id {
			return emp, nil
		}
	}
	return Employee{}, errors.New("employee not found")
}

func (s *employeeService) UpdateEmployee(id int, name, dob string) (Employee, error) {
	for i, emp := range s.employees {
		if emp.ID == id {
			s.employees[i].Name = name
			s.employees[i].DOB = dob
			return s.employees[i], nil
		}
	}
	return Employee{}, errors.New("employee not found")
}

func (s *employeeService) DeleteEmployee(id int) error {
	for i, emp := range s.employees {
		if emp.ID == id {
			s.employees = append(s.employees[:i], s.employees[i+1:]...)
			return nil
		}
	}
	return errors.New("employee not found")
}

func (s *employeeService) ListEmployees() ([]Employee, error) {
	return s.employees, nil
}
