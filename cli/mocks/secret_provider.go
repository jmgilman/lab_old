package mocks

type MockSecretProvider struct {
	FnDelete   func(key string) error
	FnGenerate func(key string) error
	FnGet      func(key string) (string, error)
	FnSet      func(key string, value string) error
}

func (m *MockSecretProvider) Delete(key string) error {
	return m.FnDelete(key)
}

func (m *MockSecretProvider) Generate(key string) error {
	return m.FnGenerate(key)
}

func (m *MockSecretProvider) Get(key string) (string, error) {
	return m.FnGet(key)
}

func (m *MockSecretProvider) Set(key string, value string) error {
	return m.FnSet(key, value)
}
