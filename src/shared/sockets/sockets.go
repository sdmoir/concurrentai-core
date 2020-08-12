package sockets

// Listener : A socket listener
type Listener interface {
	Read(socketAddress string) ([]byte, error)
}

// Writer : A socket writer
type Writer interface {
	Write(socketAddress string, data []byte) error
}
