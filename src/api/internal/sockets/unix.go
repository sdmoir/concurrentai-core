package sockets

import (
	"io/ioutil"
	"net"
)

// UnixListener : A unix socket listener
type UnixListener struct{}

// NewUnixListener : Create a new unix socket listener
func NewUnixListener() *UnixListener {
	return &UnixListener{}
}

// Read : Read data from a unix socket
func (unix *UnixListener) Read(socketAddress string) ([]byte, error) {
	listener, err := net.Listen("unix", socketAddress)
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	return data, nil
}
