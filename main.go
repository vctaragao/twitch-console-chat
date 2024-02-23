package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

const (
	TwitchOAuthUrl = "https://id.twitch.tv/oauth2/authorize"

	ClientIDKey = "CONSOLE_CHAT_CLIENT_ID"
	RedirectURI = "http://localhost:7777/redirect"
)

type (
	twitchOAuthParams struct {
		clientID     string
		forceVerify  bool
		redirectURI  string
		responseType string
		scope        string
		state        string
	}
)

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

func webAuth(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv(ClientIDKey)
	if clientID == "" {
		log.Fatal("client_id empty")
	}

	params := twitchOAuthParams{
		clientID:     clientID,
		forceVerify:  false,
		redirectURI:  RedirectURI,
		responseType: "code",
		scope:        "chat:read",
		state:        uuid.New().String(),
	}

	url, err := url.Parse(params.parseIntoUrl(TwitchOAuthUrl))
	log.Printf("Url for auth: %s", url.String())
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("xdg-open", url.String())

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/auth/web", webAuth)
	http.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Redirect request received")

		fmt.Printf("URL: %v\n", r.URL.String())

		fmt.Println("Query Values: ")
		for key, values := range r.URL.Query() {
			fmt.Printf("%s: %v\n", key, values)
		}
	})

	log.Println("Server started: listening on the port 7777... ")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatal(err)
	}
}
