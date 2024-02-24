package dbtx

import "database/sql"

type DbTx interface {
	Begin() (*sql.Tx, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Close() error
}
