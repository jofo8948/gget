package strategy

import "log"

// ToStdOut provides a strategy for printing file contents to std out
type ToStdOut struct{}

// Handle is a strategy for printing file contents out to std out
func (t *ToStdOut) Handle(b []byte) error {
	log.Printf(string(b))
	return nil
}
