package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jofo8948/gget/src"
)

const (
	Workers = 100
)

func main() {
	infile := flag.String("i", "urls.txt", "Use a line-separated file to fetch many files at once.")

	flag.Parse()

	f, err := os.Open(*infile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	lines := make(chan string, Workers*3)
	wg := new(sync.WaitGroup)
	wg.Add(Workers)

	for i := 0; i < Workers; i++ {
		go func() {
			workerFn(lines)
			wg.Done()
		}()
	}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines <- sc.Text()
	}

	close(lines)
	wg.Wait()
}

func workerFn(lines <-chan string) {
	parseURL := func(line string) (u *url.URL) {
		var err error
		u, err = url.Parse(strings.TrimSpace(line))
		if err != nil {
			log.Printf("%s is not a valid URL.", line)
			return nil
		}
		return u
	}
	for line := range lines {
		u := parseURL(line)
		if u != nil {
			download.File(u, filepath.Join(u.Host, filepath.Join(strings.Split(u.Path, "/")...)))
		}
	}
}
