package worker

type Worker struct {
	EpollFD  int
	ClientFD int
}
