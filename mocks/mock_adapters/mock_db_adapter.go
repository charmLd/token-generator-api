package mock_adapters

import (
	"context"
	"database/sql"
	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/util/config"
)

type MockAdapter struct {
	cfg  config.DBConfig
	Pool *sql.DB
}

func (a *MockAdapter) GetStats() (s sql.DBStats) {
	return
}

// New creates a new Postgres adapter instance.
func (a *MockAdapter) New(config config.DBConfig) (adapters.DBAdapterInterface, error) {
	return &MockAdapter{}, nil
}
func (a *MockAdapter) Query(query string, parameters map[string]interface{}) ([]map[string]interface{}, error) {
	return nil, nil
}
func (a *MockAdapter) RowQuery(query string, parameters map[string]interface{}) (*sql.Rows, error) {
	return nil, nil
}
func (a *MockAdapter) Exec(query string, parameters map[string]interface{}) (sql.Result, error) {
	return nil, nil
}
func (a *MockAdapter) BeginTransaction(ctx context.Context) (adapters.TransactionInterface, error) {
	return nil, nil
}
func (a *MockAdapter) Prepare(ctx context.Context, sql string) (adapters.StatementInterface, error) {
	return nil, nil
}
func (a *MockAdapter) Destruct() {
}
