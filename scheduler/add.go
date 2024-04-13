package scheduler

import "my-server/worker"

func (s *Scheduler) Add(epollFD, clientFD int) {
	//s.Workers/ <- Task{
	//	EpollFD:  epollFD,
	//	ClientFD: clientFD,
	//}
	s.Workers <- worker.Worker{
		EpollFD:  epollFD,
		ClientFD: clientFD,
	}
}
