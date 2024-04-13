package server

import (
	"fmt"
	"golang.org/x/sys/unix"
	"my-server/worker"
)

func (s *Server) Start() error {
	events := make([]unix.EpollEvent, maxEvents)
	for {
		countEvents, err := unix.EpollWait(s.EpollFD, events, -1)
		if err != nil {
			return fmt.Errorf("server start: epoll wait: %w", err)
		}

		for countEvent := 0; countEvent < countEvents; countEvent++ {
			fd := int(events[countEvent].Fd)
			if fd == s.ServerFD {
				clientFD, _, err := unix.Accept(fd)
				if err != nil {
					return fmt.Errorf("server start: accept unix socket: %w", err)
				}

				err = unix.SetNonblock(fd, true)
				if err != nil {
					return fmt.Errorf("server start: set noblock: %w", err)
				}

				err = unix.EpollCtl(s.EpollFD, unix.EPOLL_CTL_ADD, clientFD, &unix.EpollEvent{
					Events: unix.EPOLLIN | unix.EPOLLET,
					Fd:     int32(clientFD),
				})
				if err != nil {
					return fmt.Errorf("server start: epoll ctl: %w", err)
				}
			} else {
				wrk := worker.New(fd)

				err = wrk.Work()
				if err != nil {
					return fmt.Errorf("server start: work: %w", err)
				}

				err = unix.EpollCtl(s.EpollFD, unix.EPOLL_CTL_DEL, fd, nil)
				if err != nil {
					return fmt.Errorf("server start: delete client fd from epoll: %w", err)
				}

				err := unix.Close(fd)
				if err != nil {
					return fmt.Errorf("server start: close client fd: %w", err)
				}
			}
		}
	}
}
