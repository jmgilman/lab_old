package mocks

import "io"

type MockImageProvider struct {
	FnFetch    func(channel, arch, filename string) (io.ReadCloser, int64, error)
	FnValidate func(data io.ReadCloser, channel, arch, filename string) error
}

func (m *MockImageProvider) Fetch(channel, arch, filename string) (io.ReadCloser, int64, error) {
	return m.FnFetch(channel, arch, filename)
}

func (m *MockImageProvider) Validate(data io.ReadCloser, channel, arch, filename string) error {
	return m.FnValidate(data, channel, arch, filename)
}
