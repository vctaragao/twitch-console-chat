package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vctaragao/twitch-chat/internal/twauth"
	"github.com/vctaragao/twitch-chat/pkg/browser"
)

const (
	ClientSecretKey = "CONSOLE_CHAT_SECRET"
	ClientIDKey     = "CONSOLE_CHAT_CLIENT_ID"
)

func main() {
	shoutdown := make(chan os.Signal, 1)
	signal.Notify(shoutdown, os.Interrupt, syscall.SIGTERM)

	clientID := os.Getenv(ClientIDKey)
	if clientID == "" {
		log.Fatal("client_id empty")
	}

	secret := os.Getenv(ClientSecretKey)
	if secret == "" {
		log.Fatal("client_secret empty")
	}

	sqliteDB, err := twauth.NewDefaultSqliteDB()
	if err != nil {
		log.Fatalf("unable to create sqlite db: %v", err)
	}

	twitchAuthServer := twauth.NewServer(twauth.TwitchAuthParams{
		Repo:           sqliteDB,
		ClientID:       clientID,
		Secret:         secret,
		Port:           ":7777",
		BrowserHandler: browser.NewBrowserHandler(),
		Logger:         log.New(os.Stdout, "twitch-console-chat", log.LstdFlags),
	})

	twitchAuthServer.Start()

	<-shoutdown

	if err := sqliteDB.Db.Close(); err != nil {
		log.Printf("unable to close sqlite db: %v", err)
	}

	if err := twitchAuthServer.Close(); err != nil {
		log.Printf("unable to close server: %v", err)
	}
}
