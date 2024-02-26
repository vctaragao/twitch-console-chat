package integration

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vctaragao/twitch-chat/internal/twauth"
	"github.com/vctaragao/twitch-chat/internal/twauth/infra/database"
	"github.com/vctaragao/twitch-chat/pkg/dbtx"
)

const (
	ClientSecretKey = "CONSOLE_CHAT_SECRET"
	ClientIDKey     = "CONSOLE_CHAT_CLIENT_ID"
)

func (s *TestSuite) main() error {
	clientID := os.Getenv(ClientIDKey)
	if clientID == "" {
		return fmt.Errorf("client_id empty")
	}

	secret := os.Getenv(ClientSecretKey)
	if secret == "" {
		return fmt.Errorf("client_secret empty")
	}

	dbConn, sqliteDB, err := newInMemoryTestDB()
	if err != nil {
		return fmt.Errorf("unable to create sqlite db: %+v", err)
	}

	if err := sqliteDB.CreateTables(); err != nil {
		return fmt.Errorf("unable to create tables: %v", err)
	}

	browserHandler := &browserHandlerMock{}

	logBuf := &bytes.Buffer{}
	bufferLogger := log.New(logBuf, "twitch-console-chat", log.LstdFlags)
	twitchAuthServer := twauth.NewServer(twauth.TwitchAuthParams{
		Repo:           sqliteDB,
		ClientID:       clientID,
		Secret:         secret,
		Port:           ":7777",
		BrowserHandler: browserHandler,
		Logger:         bufferLogger,
	})

	twitchAuthServer.Start()

	if ready, err := serverReady(); !ready {
		return fmt.Errorf("server not ready: %w", err)
	}

	s.dbConn = dbConn
	s.sqliteDB = sqliteDB
	s.bufferLogger = logBuf
	s.twitchAuthServer = twitchAuthServer
	s.mocks["browserHandler"] = browserHandler

	return nil
}

func newInMemoryTestDB() (*sql.DB, *database.SqliteAdapter, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to open sqlite db: %w", err)
	}

	sqliteDB := database.SqliteAdapter{Db: db}

	if err := sqliteDB.CreateTables(); err != nil {
		return nil, nil, fmt.Errorf("unable to create tables: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to begin transaction: %w", err)
	}

	sqliteAdapter := &database.SqliteAdapter{Db: &dbtx.TxTestAdapter{Tx: tx}}

	return db, sqliteAdapter, nil
}

func serverReady() (bool, error) {
	ready := false
	for i := 0; i < 10; i++ {
		log.Println("Checking server status...")
		request, err := http.NewRequest("GET", "http://localhost:7777/healthz", nil)
		if err != nil {
			return ready, fmt.Errorf("unable to create request: %w", err)
		}

		client := http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			return ready, fmt.Errorf("unable to send request: %w", err)
		}
		log.Println("resp: ", resp)

		if resp.StatusCode == 200 {
			ready = true
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	if !ready {
		return ready, errors.New("server not ready after 5 seconds: timeout error")
	}

	return ready, nil
}
