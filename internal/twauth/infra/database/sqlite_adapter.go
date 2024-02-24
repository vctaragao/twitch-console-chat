package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/vctaragao/twitch-chat/internal/twauth/entity"
	"github.com/vctaragao/twitch-chat/pkg/dbtx"
)

type SqliteAdapter struct {
	Db dbtx.DbTx
}

func NewDB() (*SqliteAdapter, error) {
	db, err := sql.Open("sqlite3", "./console_chat.Db")
	if err != nil {
		return nil, err
	}

	return &SqliteAdapter{Db: db}, nil
}

func (a *SqliteAdapter) CreateTables() error {
	_, err := a.Db.Exec(`CREATE TABLE IF NOT EXISTS auths (
        id INTEGER PRIMARY KEY,
        username TEXT,
        access_token TEXT,
        refresh_token TEXT,
        expires_in INTEGER,
        scope TEXT,
        token_type TEXT
    )`)

	if err != nil {
		return fmt.Errorf("unable to create auth table: %w", err)
	}

	_, err = a.Db.Exec(`CREATE TABLE IF NOT EXISTS scopes (
        id INTEGER PRIMARY KEY,
        auth_id INTEGER,
        scope TEXT,
        FOREIGN KEY(auth_id) REFERENCES auth(id)
    )`)

	return fmt.Errorf("unable to create scopes table: %w", err)
}

func (a *SqliteAdapter) InsertAuth(auth entity.AuthToken) (int, error) {
	result, err := a.Db.Exec(
		`INSERT INTO auths (access_token, expires_in, refresh_token, token_type) VALUES (?, ?, ?, ?)`,
		auth.AccessToken,
		auth.ExpiresIn,
		auth.RefreshToken,
		auth.TokenType,
	)

	if err != nil {
		return 0, fmt.Errorf("unable to insert auth token into db: %w", err)
	}

	authID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("unable to get last insert id: %w", err)
	}

	return int(authID), nil
}

func (a *SqliteAdapter) InsertScopes(auth entity.AuthToken) error {
	dbTx, err := a.Db.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}

	for _, scope := range auth.Scopes {
		if _, err = dbTx.Exec(`INSERT INTO scopes (auth_id, scope) VALUES (?, ?)`, auth.ID, scope); err != nil {
			err = fmt.Errorf("unable to insert scope into db: %v", err)
			break
		}
	}

	if err != nil {
		if err := dbTx.Rollback(); err != nil {
			return fmt.Errorf("unable to rollback transaction: %w", err)
		}

		return err
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}
