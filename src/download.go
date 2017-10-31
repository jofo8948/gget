package download

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// File will perform an HTTP GET on the specified URL, and return the response as a []byte
// or an error
func File(uri *url.URL, dst string) {
	var (
		b   []byte
		err error
	)
	if b, err = getFile(uri); err != nil {
		log.Print(err.Error())
	}

	if err := writeFile(b, dst); err != nil {
		log.Print(err.Error())
	}
}

func getFile(url *url.URL) (b []byte, err error) {
	var resp *http.Response
	if resp, err = http.Get(url.String()); err != nil {
		err = fmt.Errorf("Error fetching %s: %v", url.String(), err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error fetching %s. Error Code was: %d", url.String(), resp.StatusCode)
	}
	if err != nil {
		return nil, err
	}

	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		panic(fmt.Errorf("%s: %s", "Critical error reading response body", err.Error()))
	}
	defer resp.Body.Close()

	return b, nil
}

func writeFile(b []byte, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0666); err != nil {
		panic(err.Error())
	}
	f, err := os.Create(dst)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Error creating file %s: %v", dst, err))
		return err
	}
	defer f.Close()
	io.Copy(f, bytes.NewBuffer(b))
	return nil
}
