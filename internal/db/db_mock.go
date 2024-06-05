package db

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func (m *MockDB) QueryRow(ctx context.Context, query string, args ...interface{}) *MockRow {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0).(*MockRow)
}

func (m *MockDB) Exec(ctx context.Context, query string, args ...interface{}) (interface{}, error) {
	callArgs := m.Called(ctx, query, args)
	return callArgs.Get(0), callArgs.Error(1)
}
