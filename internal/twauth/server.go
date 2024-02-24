package twauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type (
	TwitchAuthServer struct {
		repo     Repository
		server   *http.Server
		clientID string
		secret   string
		port     string
	}

	TwitchAuthParams struct {
		Repo     Repository
		ClientID string
		Secret   string
		Port     string
	}
)

func NewServer(params TwitchAuthParams) *TwitchAuthServer {
	return &TwitchAuthServer{
		repo:     params.Repo,
		clientID: params.ClientID,
		secret:   params.Secret,
		port:     params.Port,
	}
}

func (s *TwitchAuthServer) Start() {
	authTokenService := NewAuthTokenService(s.repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if err := Authenticate(s.clientID); err != nil {
			log.Printf("unable to authenticate user: %v\n", err)
		}
	})

	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		var code string
		for key, values := range r.URL.Query() {
			if key == "code" {
				code = values[0]
			}
		}

		if err := authTokenService.AuthToken(s.clientID, s.secret, code); err != nil {
			log.Printf("unable to fetch auth token: %v\n", err)
		}
	})

	server := http.Server{
		Addr:    s.port,
		Handler: mux,
	}

	s.server = &server

	go func() {
		log.Printf("Twitch auth server started: listening on the port %s...\n", s.port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *TwitchAuthServer) Close() error {
	if err := s.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("unable to shutdown server: %w", err)
	}

	return nil
}