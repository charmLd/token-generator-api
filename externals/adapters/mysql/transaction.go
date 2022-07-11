package mysql

import (
	"database/sql"
	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
)

// Transaction implements sql transaction interface
type Transaction struct {
	txn *sql.Tx
}

// PrepareAll prepares multiple statements
func (t *Transaction) PrepareAll(queries map[string]string) (map[string]*sql.Stmt, error) {

	prepStatements := make(map[string]*sql.Stmt)

	for key, queryString := range queries {

		prepStatement, err := t.txn.Prepare(queryString)

		if err != nil {
			return nil, err
		}

		prepStatements[key] = prepStatement
	}

	return prepStatements, nil
}

// Commit commits SQL transaction
func (t *Transaction) Commit() error {
	err := t.txn.Commit()
	if err != nil {

		return err
	}

	return nil
}

// Prepare prepares SQL statement in a transaction
func (t *Transaction) Prepare(query string) (adapters.StatementInterface, error) {

	stmt, err := t.txn.Prepare(query)

	if err != nil {
		return nil, err
	}

	return &Statement{
		stmt: stmt,
	}, nil
}

// Rollback reverts the transaction
func (t *Transaction) Rollback() error {
	return t.txn.Rollback()
}
