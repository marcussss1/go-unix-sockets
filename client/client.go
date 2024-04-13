package client

import (
	"fmt"
	"golang.org/x/sys/unix"
	"my-server/utils"
)

type Client struct {
	SocketFD   int
	ServerAddr *unix.SockaddrInet4
}

func New(port int, addr [4]byte) (*Client, error) {
	socketFD, err := unix.Socket(
		unix.AF_INET,
		unix.SOCK_STREAM,
		unix.IPPROTO_IP,
	)
	if err != nil {
		return nil, fmt.Errorf("client new: create unix socket: %w", err)
	}

	serverAddr := &unix.SockaddrInet4{
		Port: port,
		Addr: addr,
	}

	err = unix.Connect(socketFD, serverAddr)
	if err != nil {
		return nil, fmt.Errorf("client new: connect to server: %w", err)
	}

	return &Client{
		SocketFD:   socketFD,
		ServerAddr: serverAddr,
	}, nil
}

func (c *Client) Close() {
	err := utils.CloseSocket(c.SocketFD)
	if err != nil {
		fmt.Printf("client close: close socket %v\n", err)
		return
	}

	fmt.Printf("close socket %d\n", c.SocketFD)
}
