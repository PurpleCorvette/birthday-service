package employee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddEmployee(t *testing.T) {
	service := NewEmployeeService()

	employee, err := service.AddEmployee("John Doe", "1990-01-01")
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", employee.Name)

	employee, err = service.AddEmployee("Jane Doe", "1992-02-02")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", employee.Name)
}

func TestGetEmployee(t *testing.T) {
	service := NewEmployeeService()
	service.AddEmployee("John Doe", "1990-01-01")

	employee, err := service.GetEmployee(1)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", employee.Name)

	_, err = service.GetEmployee(2)
	assert.Error(t, err)
}
