package scheduler

import (
	"fmt"
	"my-server/worker"
)

func (s *Scheduler) Start() {
	for countWorker := 0; countWorker < s.CountWorkers; countWorker++ {
		go s.startWorker()
	}

	for {
		select {
		case task := <-s.Tasks:
			s.Workers <- worker.Worker{
				Task: task,
			}
		}
	}
}

func (s *Scheduler) startWorker() {
	for {
		select {
		case wrk := <-s.Workers:
			err := wrk.Work()
			if err != nil {
				fmt.Printf("server start: work: %v\n", err)
			}
		}
	}
}
