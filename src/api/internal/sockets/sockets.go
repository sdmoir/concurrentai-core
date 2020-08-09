package sockets

// Listener : A socket listener
type Listener interface {
	Read(socketAddress string) ([]byte, error)
}
