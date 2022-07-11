package mysql

import (
	"database/sql"
)

// Statement type is used to implement sql statement interface
type Statement struct {
	stmt *sql.Stmt
}

// Exec executes a prepared statement
func (st *Statement) Exec(args ...interface{}) (sql.Result, error) {


	data, err := st.stmt.Exec(args...)

	if err != nil {

		return nil, err
	}


	return data, nil
}

// Query executes a prepared statement and returns sql.Rows type
func (st *Statement) Query(args ...interface{}) (*sql.Rows, error) {

	result, err := st.stmt.Query(args...)

	if err != nil {


		return nil, err
	}


	return result, nil
}

// InsertAndGetID Execute INSERT query and returns the last inserted ID
func (st *Statement) InsertAndGetID(args ...interface{}) (int, error) {

	lastInsertedID := 0
	err := st.stmt.QueryRow(args...).Scan(&lastInsertedID)

	if err != nil {
		return 0, err
	}

	return lastInsertedID, nil

}
// Close closes the statement
func (st *Statement) Close() error {
	return st.stmt.Close()
}
