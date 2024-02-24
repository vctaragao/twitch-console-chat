package twauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/vctaragao/twitch-chat/internal/twauth/entity"
)

type AuthTokenService struct {
	repo Repository
}

func NewAuthTokenService(repo Repository) *AuthTokenService {
	return &AuthTokenService{repo: repo}
}

func (s *AuthTokenService) AuthToken(clientID, secret, code string) error {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", secret)
	data.Set("redirect_uri", RedirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, TwitchOauthTokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("unable to create token auth req: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: time.Second * 3}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to fecth auth token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var buf bytes.Buffer
		n, err := buf.ReadFrom(resp.Body)
		if err != nil {
			return fmt.Errorf("unable to read response body: %w", err)
		}

		return fmt.Errorf("unable to fetch auth token: %v, body: %v", resp.Status, buf.String()[:n])
	}

	var authToken entity.AuthToken
	if err := json.NewDecoder(resp.Body).Decode(&authToken); err != nil {
		return fmt.Errorf("unable to decode auth token response: %w", err)
	}

	authID, err := s.repo.InsertAuth(authToken)
	if err != nil {
		return fmt.Errorf("error while inserting auth token: %w", err)
	}

	authToken.ID = authID

	if err := s.repo.InsertScopes(authToken); err != nil {
		return fmt.Errorf("error while inserting scopes: %w", err)
	}

	return nil
}
