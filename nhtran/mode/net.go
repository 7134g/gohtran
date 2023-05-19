package mode

import (
	"net"
)

func createServer(address string) (net.Listener, error) {
	server, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return server, err
}

func createDial(address string) (net.Conn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return conn, err
}

func accept(listener net.Listener) (net.Conn, error) {
	conn, err := listener.Accept()
	if conn == nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return conn, nil
}
