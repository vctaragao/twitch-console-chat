package assertdb_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func (s *AssertDBSuite) TestSeedTable() {
	t := s.T()

	user, err := s.testDBHandler.SeedTable(t, "users", map[string]any{
		"id":   uuid.New().String(),
		"name": "Teste",
	})
	require.NoError(t, err)

	s.testDBHandler.AssertInTable(t, "users", map[string]any{
		"id": user["id"],
	})
}

func (s *AssertDBSuite) TestCountInTable() {
	t := s.T()

	userID := "1eee9ec6-2237-6a70-f1ac-ed2591cc8be6"

	s.testDBHandler.CountInTable(t, "users", 0, map[string]any{
		"id": userID,
	})

	_, err := s.testDBHandler.SeedTable(t, "users", map[string]any{
		"id": userID,
	})
	require.NoError(t, err)

	s.testDBHandler.AssertInTable(t, "users", map[string]any{
		"id": "1eee9ed0-559d-6b80-7a92-588444af3b04",
	})
}
