package mocks

type MockApp struct {
	Data  interface{}
	Error error
}

func (m *MockApp) Exit(data interface{}, err error) error {
	m.Data = data
	m.Error = err

	return nil
}
