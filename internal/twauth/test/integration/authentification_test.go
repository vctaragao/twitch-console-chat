package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *TestSuite) TestAuthentification(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:7777/auth", nil)
	if err != nil {
		t.Fatalf("unable to create request: %v", err)
	}

	browserHandler := s.mocks["browserHandler"].(*browserHandlerMock)

	_ = "http://localhost:7777/auth?client_id=client_id&force_verify=false&redirect_uri=http://localhost:7777/redirect&response_type=code&scope=chat:read&state=state"
	mockCall := browserHandler.
		On("Open", mock.Anything).
		Return(nil).
		Times(1)

	fmt.Println(mockCall.Arguments[0].(string))
	client := http.Client{}
	resp, err := client.Do(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockCall.Unset()
}
