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
