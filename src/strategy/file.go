package strategy

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type toFile struct {
	destination string
}

// ToFile returns a strategy for persisting a file
// it will either  write the contents of the file to the destination given,
// or if that destination does not exist, it will return an error
func ToFile(dst string) *toFile {
	return &toFile{destination: dst}
}

func (t *toFile) Handle(b []byte) (err error) {
	return writeFile(b, t.destination)
}

func writeFile(b []byte, dst string) (err error) {
	var (
		f *os.File
	)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		log.Printf(err.Error())
	}
	if f, err = os.Create(dst); err != nil {
		return fmt.Errorf(fmt.Sprintf("Error creating file %s: %v", dst, err))
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer(b))
	return err
}
