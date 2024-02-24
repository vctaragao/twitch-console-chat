package dbtx

import "database/sql"

type DbAdapter struct {
	DB *sql.DB
}

func NewDB(db *sql.DB) *DbAdapter {
	return &DbAdapter{DB: db}
}

func (a *DbAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.DB.Exec(query, args...)
}

func (a *DbAdapter) Begin() (*sql.Tx, error) {
	return a.DB.Begin()
}

func (a *DbAdapter) Close() error {
	return a.DB.Close()
}
