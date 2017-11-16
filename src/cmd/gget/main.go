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

	"github.com/jofo8948/gget/src/gget"
	"github.com/jofo8948/gget/src/strategy"
)

const workers = 100

func main() {
	var (
		f   *os.File
		err error
	)

	infile := flag.String("i", "urls.txt", "Use a line-separated file to fetch many files at once.")
	flag.Parse()

	if f, err = os.Open(*infile); err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	lines := make(chan string, workers*3)
	wg := new(sync.WaitGroup)
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			parseURL := func(line string) (u *url.URL) {
				var err error
				if line == "" {
					return nil
				}
				u, err = url.Parse(strings.TrimSpace(line))
				if err != nil {
					log.Printf("%s is not a valid URL.", line)
					return nil
				}
				return u
			}
			if url := parseURL(<-lines); url != nil {
				dst := filepath.Join(url.Host, filepath.Join(strings.Split(url.Path, "/")...))
				download(&gget.GGet{URL: url, Strategy: strategy.ToFile(dst)})
			}
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

func download(gget *gget.GGet) {
	if err := gget.Execute(); err != nil {
		log.Printf("error retrieving file from URI %v, error %s", gget.URL, err.Error())
	}
}
