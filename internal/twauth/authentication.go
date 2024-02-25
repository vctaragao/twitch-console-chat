package twauth

import (
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/vctaragao/twitch-chat/pkg/browser"
)

type twitchOAuthParams struct {
	clientID     string
	forceVerify  bool
	redirectURI  string
	responseType string
	scope        string
	state        string
}

func (t *twitchOAuthParams) parseIntoUrl(url string) string {
	return fmt.Sprintf("%s?client_id=%s&force_verify=%v&redirect_uri=%s&response_type=%s&scope=%s&state=%s",
		url,
		t.clientID,
		t.forceVerify,
		t.redirectURI,
		t.responseType,
		t.scope,
		t.state,
	)
}

func Authenticate(clientID string, browserHandler browser.BrowserHandler) error {
	params := twitchOAuthParams{
		clientID:     clientID,
		forceVerify:  false,
		redirectURI:  RedirectURI,
		responseType: "code",
		scope:        "chat:read",
		state:        uuid.New().String(),
	}

	url, err := url.Parse(params.parseIntoUrl(TwitchOAuthUrl))
	if err != nil {
		return fmt.Errorf("unable to parse url: %w", err)
	}

	if err := browserHandler.Open(url.String()); err != nil {
		return fmt.Errorf("unable to open browser for user authentification: %w", err)
	}

	return nil
}
