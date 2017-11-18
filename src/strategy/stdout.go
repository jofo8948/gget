package strategy

import "log"

// ToStdOut provides a strategy for printing file contents to std out
type ToStdOut struct {
	Size int
}

// Handle is a strategy for printing file contents out to std out
func (t *ToStdOut) Handle(b []byte) error {
	if t.Size != 0 {
		log.Printf(string(b)[:t.Size])
	} else {
		log.Printf(string(b))
	}

	return nil
}
