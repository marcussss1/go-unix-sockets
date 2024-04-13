package scheduler

import (
	"fmt"
	"my-server/worker"
	"sync"
	"time"
)

type Scheduler struct {
	Tasks   chan Task
	Workers chan worker.Worker
}

func New(countWorkers int) *Scheduler {
	workers := make(chan worker.Worker, countWorkers)
	wg := &sync.WaitGroup{}
	wg.Add(countWorkers)

	for countWorker := 0; countWorker < countWorkers; countWorker++ {
		go func() {
			wg.Done()
			for {
				select {
				case wrk := <-workers:
					err := wrk.Work()
					if err != nil {
						fmt.Printf("server start: work: %v\n", err)
					}
				default:
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}

	wg.Wait()

	return &Scheduler{
		Tasks:   make(chan Task, 1024),
		Workers: workers,
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
