package sockets

import (
	"log"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

// UnixListener : A unix socket listener
type UnixListener struct{}

// NewUnixListener : Create a new unix socket listener
func NewUnixListener() *UnixListener {
	return &UnixListener{}
}

// Read : Read data from a unix socket
func (unixListener *UnixListener) Read(socketAddress string) ([]byte, error) {
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

// UnixWriter : A unix socket writer
type UnixWriter struct {
	SocketDiscoveryIntervalMilliseconds int
	SocketDiscoveryTimeoutMilliseconds  int
}

// NewUnixWriter : Create a new unix socket writer
func NewUnixWriter() *UnixWriter {
	return &UnixWriter{
		SocketDiscoveryIntervalMilliseconds: 10,
		SocketDiscoveryTimeoutMilliseconds:  30000,
	}
}

// Write : Write data to a unix socket
func (unixWriter *UnixWriter) Write(socketAddress string, data []byte) error {
	timeout := unixWriter.SocketDiscoveryTimeoutMilliseconds
	interval := unixWriter.SocketDiscoveryIntervalMilliseconds

	// Wait for socket to exist if it does not already
	for i := 0; i <= timeout; i += interval {
		log.Println("Checking for socket " + socketAddress)
		if _, err := os.Stat(socketAddress); os.IsNotExist(err) {
			if (i + interval) < timeout {
				time.Sleep(10 * time.Millisecond)
			} else {
				return errors.Wrap(err, "timed out waiting for socket")
			}
		}
	}

	connection, err := net.Dial("unix", socketAddress)
	if err != nil {
		return errors.Wrap(err, "failed to connect to socket")
	}
	defer connection.Close()

	if _, err := connection.Write(data); err != nil {
		return errors.Wrap(err, "failed to write data to socket")
	}

	return nil
}
