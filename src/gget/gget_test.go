package gget

import (
	"net/url"
	"testing"
)

type mockGetter struct {
	mockData []byte
	err      error
}

func (mg *mockGetter) Get(u *url.URL) ([]byte, error) {
	if mg.err != nil {
		return nil, mg.err
	}
	return mg.mockData, nil
}

func testGGet(t *testing.T) {
	mock := mockGetter{mockData: []byte{0, 1, 1, 1, 1}, err: nil}
}
