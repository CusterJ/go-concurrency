package workers

import "sync"

type disp struct {
	worker
	wg            *sync.WaitGroup
	jobChannel    chan Job
	ResultChannel chan JobResult
	urls          []string
}

type Job struct {
	ID   int
	Site string
}

type JobResult struct {
	WorkerID int
	JobID    int
	Site     string
	Time     float64
	Response string
}

func NewDisp(urls []string) *disp {
	return &disp{
		urls:          urls,
		ResultChannel: make(chan JobResult),
	}
}

func (d *disp) Start(workers int) {
	d.wg = &sync.WaitGroup{}
	d.jobChannel = make(chan Job)

	for i := 1; i <= workers; i++ {
		go d.worker.Start3(i, d.wg, d.jobChannel, d.ResultChannel)
	}
	d.wg.Add(workers)

	for i, v := range d.urls {
		d.jobChannel <- Job{
			ID:   i,
			Site: v,
		}
	}
	close(d.jobChannel)
	d.wg.Wait()
	close(d.ResultChannel)
}
