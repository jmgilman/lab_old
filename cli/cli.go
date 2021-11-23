package cli

import "fmt"

type App interface {
	Exit(data interface{}, err error) error
}

type AppResult struct {
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Success bool        `json:"success"`
}

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

func Exit(message string) AppError {
	return AppError{
		err: fmt.Errorf(message),
	}
}
