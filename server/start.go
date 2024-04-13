package server

import (
	"fmt"
	"golang.org/x/sys/unix"
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
				s.Scheduler.Add(s.EpollFD, fd)
			}
		}
	}
}
