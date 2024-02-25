package integration

import "github.com/stretchr/testify/mock"

type browserHandlerMock struct {
	mock.Mock
}

func (b *browserHandlerMock) Open(url string) error {
	args := b.Called(url)
	return args.Error(0)
}
