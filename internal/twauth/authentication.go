package twauth

import (
	"fmt"
	"net/url"
	"os/exec"

	"github.com/google/uuid"
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

func Authenticate(clientID string) error {
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
		return fmt.Errorf("unable to parse url: %v", err)
	}

	if err := exec.Command("xdg-open", url.String()).Run(); err != nil {
		return fmt.Errorf("unable to open browser for user authentification: %v", err)
	}

	return nil
}
