package dbtx

import "database/sql"

type TxTestAdapter struct {
	Tx *sql.Tx
}

func NewTxTest(tx *sql.Tx) *TxTestAdapter {
	return &TxTestAdapter{Tx: tx}
}

func (a *TxTestAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return a.Tx.Exec(query, args...)
}

func (a *TxTestAdapter) Begin() (*sql.Tx, error) {
	return a.Tx, nil
}

func (a *TxTestAdapter) Commit() error {
	return nil
}

func (a *TxTestAdapter) Rollback() error {
	return nil
}

func (a *TxTestAdapter) CommitTest() error {
	return a.Tx.Commit()
}

func (a *TxTestAdapter) RollbackTest() error {
	return a.Tx.Rollback()
}

func (a *TxTestAdapter) Close() error {
	return a.Tx.Rollback()
}
