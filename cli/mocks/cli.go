package mocks

import "github.com/urfave/cli/v2"

type MockApp struct {
	Data  interface{}
	Error error
}

func (m *MockApp) Exit(c *cli.Context, data interface{}, err error) error {
	m.Data = data
	m.Error = err

	return nil
}
