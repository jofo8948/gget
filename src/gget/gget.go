package gget

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jofo8948/gget/src/strategy"
)

// retriever provides the interface for Retrieving a file
type retriever interface {
	// get will retrieve the file from the specified URL as a []byte or an error
	get(u *url.URL) (b []byte, err error)
}

type httpRetriever struct{}

// get will perform an HTTP GET on the specified URL and return the raw bytes
// this is an internal method which GGet relies on to download a file,
// have another retrieval strategy? make a PR :)
// The consumers of GGet shouldn't concern themselves with how it retrieves
// unless they want to dig into the code
func (r *httpRetriever) get(uri *url.URL) (b []byte, err error) {
	return r.getFile(uri)
}

func (r *httpRetriever) getFile(url *url.URL) (b []byte, err error) {
	var resp *http.Response
	if url == nil {
		return nil, fmt.Errorf("empty url, cannot retrieve a ghost resource")
	}
	if resp, err = http.Get(url.String()); err != nil {
		return nil, fmt.Errorf("Error fetching %s: %v", url.String(), err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching %s. Error Code was: %d", url.String(), resp.StatusCode)
	}
	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(fmt.Errorf("%s: %s", "Critical error reading response body", err.Error()))
	}
	defer resp.Body.Close()
	return b, nil
}

// GGet provides a way to retrieve the Data at URL
// and a convenient interface for persisting the data retrieved
type GGet struct {
	URL      *url.URL
	Strategy strategy.Handler

	r retriever
}

// Default returns an opinionated version of GGet
// utilizing a retriever that defaults to an HTTP GET
func Default(u *url.URL, s strategy.Handler) *GGet {
	return &GGet{URL: u, Strategy: s, r: &httpRetriever{}}
}

// Execute will execute gget with the provided strategy
func (g *GGet) Execute() (err error) {
	var b []byte

	if b, err = g.r.get(g.URL); err != nil {
		return err
	}

	return g.Strategy.Handle(b)
}
