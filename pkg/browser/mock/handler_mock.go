package mock

import "github.com/stretchr/testify/mock"

type BrowserHandlerMock struct {
	mock.Mock
}

func NewBrowserHandlerMock() *BrowserHandlerMock {
	return &BrowserHandlerMock{}
}

func (b *BrowserHandlerMock) Open(url string) error {
	args := b.Called(url)
	return args.Error(0)
}
