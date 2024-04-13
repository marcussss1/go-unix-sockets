package worker

type Worker struct {
	FD int
}

func New(fd int) Worker {
	return Worker{
		FD: fd,
	}
}
