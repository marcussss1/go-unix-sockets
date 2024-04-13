package scheduler

type Task struct {
	EpollFD  int
	ClientFD int
}
