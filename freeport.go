package spaproxy

import (
	"fmt"
	"net"
)

// GetFreePort checks and returns free port if specified port is privildeged port or zero
func GetFreePort(port int) (int, error) {
	if port < 1024 {
		port = 0
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
