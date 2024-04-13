package server

import (
	"fmt"
	"golang.org/x/sys/unix"
	"my-server/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	ServerFD        int
	EpollFD         int
	ClientSocketFDs map[int]int
	Signals         chan os.Signal
	//Scheduler       *scheduler.Scheduler
}

func New(port int, addr [4]byte) (*Server, error) {
	serverFD, err := newServerSocket(port, addr)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	epollFD, err := newEpollSocket(serverFD)
	if err != nil {
		return nil, fmt.Errorf("epoll new: epoll create: %w", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	srv := &Server{
		ServerFD:        serverFD,
		EpollFD:         epollFD,
		ClientSocketFDs: make(map[int]int),
		Signals:         signals,
		//Scheduler:       nil,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go catchSignal(wg, srv)
	wg.Wait()

	return srv, nil
}

func newServerSocket(port int, addr [4]byte) (int, error) {
	serverFD, err := unix.Socket(
		unix.AF_INET,
		unix.O_NONBLOCK|unix.SOCK_STREAM,
		unix.IPPROTO_IP,
	)
	if err != nil {
		return 0, fmt.Errorf("server new: create unix socket: %w", err)
	}

	err = unix.SetsockoptInt(serverFD, unix.SOL_SOCKET, unix.SO_REUSEPORT, port)
	if err != nil {
		return 0, fmt.Errorf("server new: set socket option: %w", err)
	}

	err = unix.SetsockoptInt(serverFD, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		return 0, fmt.Errorf("server new: set socket option: %w", err)
	}

	err = unix.SetNonblock(serverFD, true)
	if err != nil {
		return 0, fmt.Errorf("server new:set noblock: %w", err)
	}

	err = unix.Bind(serverFD, &unix.SockaddrInet4{
		Port: port,
		Addr: addr,
	})
	if err != nil {
		return 0, fmt.Errorf("server listen: bind unix socket: %w", err)
	}

	err = unix.Listen(serverFD, maxConnects)
	if err != nil {
		return 0, fmt.Errorf("server listen: listen unix socket: %w", err)
	}

	return serverFD, nil
}

func newEpollSocket(serverFD int) (int, error) {
	epollFD, err := unix.EpollCreate1(0)
	if err != nil {
		return 0, fmt.Errorf("epoll new: epoll create: %w", err)
	}

	err = unix.EpollCtl(epollFD, unix.EPOLL_CTL_ADD, serverFD, &unix.EpollEvent{
		Events: unix.EPOLLIN,
		Fd:     int32(serverFD),
	})
	if err != nil {
		return 0, fmt.Errorf("epoll new: epoll ctl: %w", err)
	}

	return epollFD, nil
}

func catchSignal(wg *sync.WaitGroup, srv *Server) {
	wg.Done()
	<-srv.Signals

	err := srv.close()
	if err != nil {
		fmt.Printf("server catch signal: close server: %v\n", err)
		return
	}

	fmt.Println("\nsucceeded to close server")
	os.Exit(0)
}

func (s *Server) close() error {
	for _, socketFD := range s.ClientSocketFDs {
		err := utils.CloseSocket(socketFD)
		if err != nil {
			return fmt.Errorf("server close: close client fds: %w", err)
		}
	}

	err := utils.CloseSocket(s.ServerFD)
	if err != nil {
		return fmt.Errorf("server close: close server fd: %w", err)
	}

	err = utils.CloseSocket(s.EpollFD)
	if err != nil {
		return fmt.Errorf("server close: close epoll fd: %w", err)
	}

	return nil
}
