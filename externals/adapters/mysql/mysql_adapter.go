package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/util/config"
	_ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
	cfg  config.DBConfig
	Pool *sql.DB
}

// New creates a new Postgres adapter instance.
func (a *Adapter) New(config config.DBConfig) (adapters.DBAdapterInterface, error) {

	a.cfg = config

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		a.cfg.User, a.cfg.Password, a.cfg.Host, a.cfg.Port, a.cfg.Database)

	db, err := sql.Open("mysql", connString)

	if err != nil {
		return nil, err
	}

	// Pool configurations
	db.SetMaxOpenConns(a.cfg.MaxOpenCon)
	db.SetMaxIdleConns(a.cfg.MaxIdleCon)

	a.Pool = db

	return a, nil
}

// Query runs a query and returns the result.
func (a *Adapter) Query(query string, parameters map[string]interface{}) ([]map[string]interface{}, error) {

	statement, err := a.Pool.Prepare(query)
	if err != nil {
		return nil, err
	}
	// check whether the query is a select statement
	if strings.ToLower(query[:1]) == "s" {

		rows, err := statement.Query(parameters)
		if err != nil {
			return nil, err
		}

		return a.prepareDataSet(rows)
	}

	result, err := statement.Exec(parameters)
	if err != nil {
		return nil, err
	}

	return a.prepareResultSet(result)
}

// RowQuery runs a query and returns the result.
func (a *Adapter) RowQuery(query string, parameters map[string]interface{}) (*sql.Rows, error) {

	statement, err := a.Pool.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(parameters)
	if err != nil {

		return nil, err
	}

	return rows, nil

}

// Exec executs a query and return a result interface.
func (a *Adapter) Exec(query string, parameters map[string]interface{}) (sql.Result, error) {
	result, err := a.Pool.Exec(query, parameters)

	if err != nil {

		return nil, err
	}

	return result, nil
}

// BeginTransaction starts a transaction
func (a *Adapter) BeginTransaction(ctx context.Context) (adapters.TransactionInterface, error) {

	tx, err := a.Pool.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}
	return &Transaction{
		txn: tx,
	}, nil
}

// Prepare prepares a row sql statement
func (a *Adapter) Prepare(ctx context.Context, sql string) (adapters.StatementInterface, error) {

	statement, err := a.Pool.PrepareContext(ctx, sql)

	if err != nil {
		return nil, err
	}

	return &Statement{
		stmt: statement,
	}, nil
}

// Destruct will close the Postgres adapter releasing all resources.
func (a *Adapter) Destruct() {
	a.Pool.Close()
}

// Prepare the resultset for all other queries.
func (a *Adapter) prepareResultSet(result sql.Result) ([]map[string]interface{}, error) {

	var data []map[string]interface{}

	row := make(map[string]interface{})

	row["affected_rows"], _ = result.RowsAffected()
	row["last_insert_id"], _ = result.LastInsertId()

	return append(data, row), nil
}

// Prepare the return dataset for select statements.
// Source: https://kylewbanks.com/blog/query-result-to-map-in-golang
func (a *Adapter) prepareDataSet(rows *sql.Rows) ([]map[string]interface{}, error) {

	defer rows.Close()

	var data []map[string]interface{}
	cols, _ := rows.Columns()

	// create a slice of interface{}'s to represent each column
	// and a second slice to contain pointers to each item in the columns slice
	columns := make([]interface{}, len(cols))
	columnPointers := make([]interface{}, len(cols))

	for i := range columns {
		columnPointers[i] = &columns[i]
	}

	for rows.Next() {
		// scan the result into the column pointers
		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		// create our map, and retrieve the value for each column from the pointers slice
		// storing it in the map with the name of the column as the key
		row := make(map[string]interface{})

		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			row[colName] = *val
		}

		data = append(data, row)
	}

	return data, nil
}
