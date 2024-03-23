package assertdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestDBHandler struct {
	ctx    context.Context
	dbConn *sql.DB
	dbTx   *sql.Tx
}

func NewTestDBHandler(ctx context.Context, dbConn *sql.DB, dbTx *sql.Tx) *TestDBHandler {
	return &TestDBHandler{
		ctx:    ctx,
		dbTx:   dbTx,
		dbConn: dbConn,
	}
}

func (h *TestDBHandler) SeedTable(t *testing.T, table string, data map[string]any) (map[string]any, error) {
	t.Helper()

	require.Greater(t, len(data), 0)

	values := make([]any, 0, len(data))
	query := fmt.Sprintf("INSERT INTO %s (", table)

	for field, value := range data {
		query += field + ","
		values = append(values, value)
	}

	query = query[:len(query)-1] + ")"

	mask := strings.Repeat("?,", len(values))
	mask = mask[:len(mask)-1]

	query += fmt.Sprintf(" VALUES (%s)", mask)

	result, err := h.dbTx.Exec(query, values...)
	require.NoError(t, err)

	if _, exists := data["id"]; !exists {
		id, err := result.LastInsertId()
		require.NoError(t, err)

		data["id"] = id
	}

	return data, nil
}

func (h *TestDBHandler) AssertInTable(t *testing.T, table string, expectedFields map[string]any) map[string]any {
	t.Helper()

	rows, err := h.dbTx.Query(h.buildWhere(t, table, expectedFields))
	require.NoError(t, err)

	columns, err := rows.Columns()
	require.NoError(t, err)

	values := make([]any, len(columns))
	for i := range values {
		var v any
		values[i] = &v
	}

	assert.Truef(t, rows.Next(), "No record was found")
	require.NoError(t, rows.Scan(values...))

	data := make(map[string]any, len(columns))
	for i, colName := range columns {
		data[colName] = *(values[i].(*any))
	}

	return data
}

func (h *TestDBHandler) buildWhere(t *testing.T, table string, fields map[string]any) string {
	t.Helper()

	query := "SELECT * FROM " + table + " WHERE "

	for field, value := range fields {
		if value == nil {
			query += field + " IS NULL"
		} else {
			query += fmt.Sprintf("%s = '%v'", field, value)
		}

		query += " AND "
	}

	query, _ = strings.CutSuffix(query, " AND ")

	return query
}

func (h *TestDBHandler) Close() error {
	return h.dbConn.Close()
}

func (h *TestDBHandler) Roolback() error {
	return h.dbTx.Rollback()
}

func (h *TestDBHandler) Begin() error {
	dbTx, err := h.dbConn.Begin()
	if err != nil {
		return err
	}

	h.dbTx = dbTx

	return nil
}

func (h *TestDBHandler) DBConn() *sql.DB {
	return h.dbConn
}
