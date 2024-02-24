package twauth

import (
	"fmt"

	"github.com/vctaragao/twitch-chat/internal/twauth/infra/database"
)

func NewDefaultSqliteDB() (*database.SqliteAdapter, error) {
	sqliteDB, err := database.NewDB()
	if err != nil {
		return nil, fmt.Errorf("unable to create sqlite db: %w", err)
	}

	if err := sqliteDB.CreateTables(); err != nil {
		return nil, fmt.Errorf("unable to create tables: %w", err)
	}

	return sqliteDB, nil
}
