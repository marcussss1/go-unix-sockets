package scheduler

import "my-server/task"

func (s *Scheduler) Add(epollFD, clientFD int) {
	s.Tasks <- task.Task{
		EpollFD:  epollFD,
		ClientFD: clientFD,
	}
}
