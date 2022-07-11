package adapters

import (
	"context"
	"database/sql"
	"github.com/charmLd/token-generator-api/util/config"
)

// DBAdapterInterface is implemented by all database adapters.
type DBAdapterInterface interface {
	// New creates a new instance of database adapter implementation.
	New(config config.DBConfig) (DBAdapterInterface, error)

	// Query runs a query and return the result map
	Query(query string, parameters map[string]interface{}) ([]map[string]interface{}, error)

	// Exec executes a query with named parameters
	Exec(query string, parameters map[string]interface{}) (sql.Result, error)

	//RowQuery runs a query and returns sql.Rows
	RowQuery(query string, parameters map[string]interface{}) (*sql.Rows, error)

	// Destruct will close the database adapter releasing all resources.
	Destruct()

	//Prepare prepares a SQL statement
	Prepare(ctx context.Context, sql string) (StatementInterface, error)

	//BeginTransaction returns a transaction interface that can be used to execute multiple SQL statements
	BeginTransaction(context.Context) (TransactionInterface, error)

	GetStats() sql.DBStats
}

// TransactionInterface will be used to execute multiple SQL statements in a reversable way
type TransactionInterface interface {

	// Prepare prepares a single SQL statement
	Prepare(query string) (StatementInterface, error)

	// PrepareAll prepares multiple sql statements at once
	PrepareAll(queries map[string]string) (map[string]*sql.Stmt, error)

	// Commit the transaction
	Commit() error

	// Rollback reverts the changes
	Rollback() error
}

// StatementInterface will be used to generalize SQL statements throughout the application
type StatementInterface interface {

	// Executed and SQL statement
	Exec(args ...interface{}) (sql.Result, error)

	Query(args ...interface{}) (*sql.Rows, error)

	// Execute an INSERT query and returns the inserted ID
	InsertAndGetID(args ...interface{}) (int, error)

	// Close the statement
	Close() error
}
