package seeder

type URLProvider interface {
	GetURL(venName string) string
}

type ActualURLProvider struct{}

func NewActualURLProvider(url string) *ActualURLProvider {
	return &ActualURLProvider{}
}

func (a *ActualURLProvider) GetURL(venName string) string {
	return "https://dgvkh-" + venName + ".nrp-nautilus.io" // Your actual URL
}

type MockURLProvider struct {
	url string
}

func NewMockURLProvider(url string) *MockURLProvider {
	return &MockURLProvider{url: url}
}

func (m *MockURLProvider) GetURL(venName string) string {
	return m.url // Your mock server URL
}
