package worker

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func (w *Worker) Work() error {
	buf := make([]byte, 1024)
	_, err := unix.Read(w.FD, buf)
	if err != nil {
		return fmt.Errorf("worker work: read from client fd:%d %v\n", w.FD, err)
	}

	var resp string
	switch string(buf[0]) {
	case "1":
		resp = "first response"
	case "2":
		resp = "second response"
	case "3":
		resp = "third response"
	default:
		resp = "unknown response"
	}

	err = unix.Send(
		w.FD,
		[]byte(resp),
		unix.MSG_DONTWAIT,
	)
	if err != nil {
		return fmt.Errorf("worker work: send message to socket: %v\n", err)
	}

	return nil
}
