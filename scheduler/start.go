package scheduler

//func (s *Scheduler) Start() {
//	for {
//		select {
//		case task := <-s.Tasks:
//			s.Workers <- worker.Worker{
//				EpollFD:  task.EpollFD,
//				ClientFD: task.ClientFD,
//			}
//		default:
//			time.Sleep(100 * time.Millisecond)
//		}
//	}
//}
