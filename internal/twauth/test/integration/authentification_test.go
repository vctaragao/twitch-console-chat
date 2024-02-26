package integration

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (s *TestSuite) TestAuthentification(t *testing.T) {
	browserHandler := s.mocks["browserHandler"].(*browserHandlerMock)

	mockCall := browserHandler.
		On("Open", mock.Anything).
		Run(func(args mock.Arguments) {
			argUrl, err := url.Parse(args.Get(0).(string))
			assert.NoError(t, err)

			assert.Equal(t, "chat:read", argUrl.Query().Get("scope"))
			assert.Equal(t, "false", argUrl.Query().Get("force_verify"))
			assert.Equal(t, "code", argUrl.Query().Get("response_type"))
			assert.Equal(t, "client_id", argUrl.Query().Get("client_id"))
			assert.NoError(t, uuid.Validate(argUrl.Query().Get("state")))
			assert.Equal(t, "http://localhost:7777/redirect", argUrl.Query().Get("redirect_uri"))
		}).
		Return(nil).
		Times(1)

	req, err := http.NewRequest("GET", "http://localhost:7777/auth", nil)
	assert.NoError(t, err)

	client := http.Client{}

	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	browserHandler.AssertExpectations(t)
	mockCall.Unset()
}
