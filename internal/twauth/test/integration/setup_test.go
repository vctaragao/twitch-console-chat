package integration

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/vctaragao/twitch-chat/internal/twauth"
	"github.com/vctaragao/twitch-chat/internal/twauth/infra/database"
)

type TestSuite struct {
	dbConn           *sql.DB
	mocks            map[string]any
	sqliteDB         *database.SqliteAdapter
	twitchAuthServer *twauth.TwitchAuthServer
}

func (s *TestSuite) beforeAll() {
	os.Setenv("CONSOLE_CHAT_SECRET", "secret")
	os.Setenv("CONSOLE_CHAT_CLIENT_ID", "client_id")

	if err := s.main(); err != nil {
		log.Fatalf("unable to setup test: %v", err)
	}
}

func TestMain(t *testing.T) {
	setup := TestSuite{mocks: make(map[string]any)}

	setup.beforeAll()

	t.Run("TestAuthentification", setup.TestAuthentification)

	setup.afterAll()
}

func (s *TestSuite) afterAll() {
	if err := s.sqliteDB.Db.Close(); err != nil {
		log.Printf("unable to final roolback test sqlite tx: %v", err)
	}

	if err := s.dbConn.Close(); err != nil {
		log.Printf("unable to close sqlite db: %v", err)
	}

	if err := s.twitchAuthServer.Close(); err != nil {
		log.Printf("unable to close server: %v", err)
	}
}
