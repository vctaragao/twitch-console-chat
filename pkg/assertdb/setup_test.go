package assertdb_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"github.com/vctaragao/twitch-chat/pkg/assertdb"
)

type AssertDBSuite struct {
	suite.Suite

	testDBHandler *assertdb.TestDBHandler
}

func (s *AssertDBSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE users (id TEXT PRIMARY KEY, name TEXT)")
	if err != nil {
		panic(err)
	}

	dbTx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	s.testDBHandler = assertdb.NewTestDBHandler(context.Background(), db, dbTx)
}

func (s *AssertDBSuite) SetupTest() {

}

func (s *AssertDBSuite) TearDownTest() {
	if err := s.testDBHandler.Roolback(); err != nil {
		panic(err)
	}

	if err := s.testDBHandler.Begin(); err != nil {
		panic(err)
	}
}

func (s *AssertDBSuite) TearDownSuite() {
	if err := s.testDBHandler.Close(); err != nil {
		panic(err)
	}
}

func TestRunSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test suite")
	}

	suite.Run(t, new(AssertDBSuite))
}
