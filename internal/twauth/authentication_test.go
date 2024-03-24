package twauth_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vctaragao/twitch-chat/internal/twauth"
	browserMock "github.com/vctaragao/twitch-chat/pkg/browser/mock"
)

func TestAuthenticate(t *testing.T) {
	cliendID := "client_id"
	browserHandler := browserMock.NewBrowserHandlerMock()

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

	assert.NoError(t, twauth.Authenticate(cliendID, browserHandler))
	mockCall.Unset()
}

func TestAuthenticateError(t *testing.T) {
	cliendID := "client_id"
	browserHandler := browserMock.NewBrowserHandlerMock()

	expectedErr := errors.New("failed")
	mockCall := browserHandler.
		On("Open", mock.Anything).
		Return(expectedErr).
		Times(1)

	err := twauth.Authenticate(cliendID, browserHandler)
	require.Error(t, err)

	openingBrowserError := &twauth.OpeningBrowserError{}
	assert.True(t, errors.As(err, openingBrowserError))

	assert.ErrorIs(t, openingBrowserError.Err, expectedErr)
	mockCall.Unset()
}
