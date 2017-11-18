package gget

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/jofo8948/gget/src/strategy"
)

const (
	PASS = true
	FAIL = false
)

func pass(msg string) bool {
	return PASS
}
func fail(msg string) bool {
	return FAIL
}

type mockRetriever struct {
	mockData []byte
	mockErr  error
}

var mockURL, _ = url.Parse("http://localhost:8080.com")

var ggetExpectError = GGet{URL: nil, Strategy: &strategy.ToStdOut{}, r: &mockRetriever{mockData: nil, mockErr: errors.New("an error")}}
var ggetExpectData = GGet{URL: mockURL, Strategy: &strategy.ToStdOut{}, r: &mockRetriever{mockData: []byte{0, 1, 1, 1}}}

func (r *mockRetriever) get(u *url.URL) ([]byte, error) {
	if r.mockErr != nil {
		return nil, r.mockErr
	}
	return r.mockData, nil
}

func TestGGet(t *testing.T) {
	testGGetExpectError := func() bool {
		err := ggetExpectError.Execute()
		if err == nil {
			return fail("testGGetExpectError")
		}

		return pass("testGGetExpectError")
	}
	testGGetExpectNoError := func() bool {
		if err := ggetExpectData.Execute(); err != nil {
			fmt.Println(err.Error())
			return fail("testGGetExpectNoError")
		}
		return pass("testGGetExpectNoError")
	}

	t1Status := testGGetExpectError()
	checkTestStatus(t1Status, t)
	t2Status := testGGetExpectNoError()
	checkTestStatus(t2Status, t)
}

func TestHttpRetriever(t *testing.T) {
	expectedResponse := []byte{0, 1, 1}
	tstServer := start([]byte{0, 1, 1})
	defer tstServer.Close()

	localhostURL, _ := url.Parse("http://localhost:8080/test")
	if response, err := new(httpRetriever).get(localhostURL); err != nil {
		t.Fail()
	} else if !bytes.Equal(response, expectedResponse) {
		t.Fail()
	}
}

func checkTestStatus(status bool, t *testing.T) {
	if status != PASS {
		t.Fail()
	}
}

func start(expectedResponse []byte) *http.Server {
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write(expectedResponse)
	})
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Server errord out: %s", err)
		}
	}()

	return server
}
