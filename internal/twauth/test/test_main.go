package test

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vctaragao/twitch-chat/internal/twauth"
	"github.com/vctaragao/twitch-chat/internal/twauth/infra/database"
	"github.com/vctaragao/twitch-chat/pkg/dbtx"
)

const (
	ClientSecretKey = "CONSOLE_CHAT_SECRET"
	ClientIDKey     = "CONSOLE_CHAT_CLIENT_ID"
)

func Main() {
	shoutdown := make(chan os.Signal, 1)
	signal.Notify(shoutdown, os.Interrupt, syscall.SIGTERM)

	clientID := os.Getenv(ClientIDKey)
	if clientID == "" {
		log.Fatal("client_id empty")
	}

	secret := os.Getenv(ClientSecretKey)
	if secret == "" {
		log.Fatal("client_secret empty")
	}

	sqliteDB, err := newInMemoryTestDB()
	if err != nil {
		log.Fatalf("unable to create sqlite db: %v", err)
	}

	if err := sqliteDB.CreateTables(); err != nil {
		log.Fatalf("unable to create tables: %v", err)
	}

	twitchAuthServer := twauth.NewServer(twauth.TwitchAuthParams{
		Repo:     sqliteDB,
		ClientID: clientID,
		Secret:   secret,
		Port:     ":7777",
	})

	twitchAuthServer.Start()

	<-shoutdown

	if err := sqliteDB.Db.Close(); err != nil {
		log.Printf("unable to close sqlite db: %v", err)
	}

	if err := twitchAuthServer.Close(); err != nil {
		log.Printf("unable to close server: %v", err)
	}
}

func newInMemoryTestDB() (*database.SqliteAdapter, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	sqliteDB := database.SqliteAdapter{Db: db}

	if err := sqliteDB.CreateTables(); err != nil {
		log.Fatalf("unable to create tables: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("unable to begin transaction: %v", err)
	}

	return &database.SqliteAdapter{Db: &dbtx.TxTestAdapter{Tx: tx}}, nil
}
