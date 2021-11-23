package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// App represents a CLI application.
type App interface {
	Exit(c *cli.Context, data interface{}, err error) error
}

// AppResult is the result returned from running a CLI command.
type AppResult struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Success bool        `json:"success"`
}

// AppError is an error returned from running a CLI command.
type AppError struct {
	err error
}

func (a AppError) Error() string {
	return a.err.Error()
}

func (a AppError) ExitCode() int {
	return 1
}

func (a AppError) Format(s fmt.State, verb rune) {
	s.Write([]byte(fmt.Sprintf("An error occurred: %s", a.err)))
}
