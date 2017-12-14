package strategy

// Handler defines an interface for the side effect that should happen upon retrieving a file
type Handler interface {
	Handle(b []byte) error
}
