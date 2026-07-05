package dispatcher

type Job struct {
	Queue string
	ID    string
}

func Next(highPriority, lowPriority <-chan Job) Job {
	select {
	case job := <-highPriority:
		return job
	case job := <-lowPriority:
		return job
	}
}
