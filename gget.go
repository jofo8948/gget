package main

import (
	"flag"
	"bufio"
	"os"
	"log"
	"sync"
	"net/url"
	"path/filepath"
	"strings"
	"net/http"
	"io"
)

const (
	Workers = 100
)

func main() {
	infile := flag.String("i", "urls.txt", "Use a line-separated file to fetch many files at once.")

	flag.Parse()

	f, err := os.Open(*infile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close();

	lines := make(chan string, Workers*3)
	wg := new(sync.WaitGroup)
	wg.Add(Workers);

	for i:=0; i < Workers; i++ {	
		go func() {
			workerFn(lines);
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(f);
	for sc.Scan() {
		lines <- sc.Text()
	}

	close(lines)
	wg.Wait()
}

func workerFn(lines <-chan string) {
	for line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			downloadFile(line)
		}
	}
}

func downloadFile(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		log.Printf("%s is not a valid URL.", uri)
		return err
	}
	path := u.Path;
	fp := filepath.Join(u.Host, filepath.Join(strings.Split(path, "/")...))

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("Error fetching %s: %v", u.String(), err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s. Error Code was: %d", u.String(), resp.StatusCode)
		return err
	}

	defer resp.Body.Close()

	os.MkdirAll(filepath.Dir(fp), 0666)

	f, err := os.Create(fp)
	if err != nil {
		log.Printf("Error creating file %s: %v", fp, err)
		return err
	}
	defer f.Close()
	io.Copy(f, resp.Body)
	return nil
}