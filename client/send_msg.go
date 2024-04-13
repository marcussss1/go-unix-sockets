package client

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func (c *Client) SendMsg(msg string) (string, error) {
	err := unix.Sendmsg(
		c.SocketFD,
		[]byte(msg),
		nil,
		c.ServerAddr,
		unix.MSG_DONTWAIT,
	)
	if err != nil {
		return "", fmt.Errorf("client send msg: send message: %w", err)
	}

	buf := make([]byte, 1024)
	_, _, err = unix.Recvfrom(
		c.SocketFD,
		buf,
		0,
	)
	if err != nil {
		return "", fmt.Errorf("client send msg: recv message: %w", err)
	}

	return string(buf), nil
}
