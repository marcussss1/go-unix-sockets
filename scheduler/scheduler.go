package scheduler

import (
	"my-server/worker"
	"time"
)

type Scheduler struct {
	//TasksCount int
	Tasks chan int
	//Workers    chan int
}

func New(countWorkers int) *Scheduler {
	tasks := make(chan int, countWorkers)
	//tasks := make(map[int]struct{})
	//workers := make(chan int, countWorkers)
	//wg := &sync.WaitGroup{}
	//wg.Add(countWorkers)

	for countWorker := 0; countWorker < countWorkers; countWorker++ {
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			for {
				select {
				case <-ticker.C:
					time.Sleep(100 * time.Millisecond)
				case task := <-tasks:
					wrk := worker.New(task)
					wrk.Work()
				}
			}
			////wg.Done()
			//for {
			//	if len(tasks) {
			//
			//	}
			//	wrk := worker.New(task)
			//	wrk.Work()
			//	time.Sleep(100 * time.Millisecond)
			//	//select {
			//	//case task := <-workers:
			//	//	wrk := worker.New(task)
			//	//	wrk.Work()
			//	//	delete(tasks, task)
			//	//}
			//}
		}()
	}

	time.Sleep(1 * time.Second)
	//wg.Wait()

	return &Scheduler{
		//TasksCount: 0,
		Tasks: tasks,
		//Workers:    workers,
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
