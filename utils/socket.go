package utils

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func CloseSocket(socketFD int) error {
	err := unix.Shutdown(socketFD, unix.SHUT_RDWR)
	if err != nil {
		return fmt.Errorf("close unix socket: %w", err)
	}

	err = unix.Close(socketFD)
	if err != nil {
		return fmt.Errorf("close unix socket: %w", err)
	}

	return nil
}
