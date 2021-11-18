package mocks

import "io"

type MockImageProvider struct {
	FnFetch    func(channel string, arch string) (io.ReadCloser, int64, error)
	FnValidate func(data io.ReadCloser, channel string, arch string) error
}

func (m *MockImageProvider) Fetch(channel string, arch string) (io.ReadCloser, int64, error) {
	return m.FnFetch(channel, arch)
}

func (m *MockImageProvider) Validate(data io.ReadCloser, channel string, arch string) error {
	return m.FnValidate(data, channel, arch)
}
