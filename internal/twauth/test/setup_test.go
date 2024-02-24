package test

import (
	"os"
	"testing"
)

func beforeAll() {
	os.Setenv("CONSOLE_CHAT_SECRET", "secret")
	os.Setenv("CONSOLE_CHAT_CLIENT_ID", "client_id")

	Main()
}

func TestMain(m *testing.M) {
	beforeAll()
	os.Exit(m.Run())
}
