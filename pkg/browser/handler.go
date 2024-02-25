package browser

import (
	"fmt"
	"os/exec"
)

type (
	BrowserHandler interface {
		Open(url string) error
	}

	BrowserOpener struct{}
)

func NewBrowserHandler() *BrowserOpener {
	return &BrowserOpener{}
}

func (b *BrowserOpener) Open(url string) error {
	if err := exec.Command("xdg-open", url).Run(); err != nil {
		return fmt.Errorf("unable to open browser, url: %v, err: %w", url, err)
	}

	return nil
}
