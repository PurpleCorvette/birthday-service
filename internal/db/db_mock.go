package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	ret := m.Called(ctx, sql, args)
	return ret.Get(0).(pgx.Row)
}
