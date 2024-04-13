package scheduler

func (s *Scheduler) Add(fd int) {
	s.Tasks <- fd
}
