package scheduler

import (
	"my-server/task"
	"my-server/worker"
)

type Scheduler struct {
	CountWorkers int
	Tasks        chan task.Task
	Workers      chan worker.Worker
}

func New(countWorkers int) *Scheduler {
	return &Scheduler{
		CountWorkers: countWorkers,
		Tasks:        make(chan task.Task, 1024),
		Workers:      make(chan worker.Worker, countWorkers),
	}
}

func (s *Scheduler) Close() {
	// todo
	//err := unix.Close(e.EpollFD)
	//if err != nil {
	//	fmt.Printf("epoll close: epoll close: %v\n", err)
	//	return
	//}
	//
	//fmt.Println("epoll close")
}
