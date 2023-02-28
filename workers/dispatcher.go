package workers

import "sync"

type disp struct {
	worker
	NumberOfWorkers int
	wg              *sync.WaitGroup
	jobChannel      chan Job
	ResultChannel   chan JobResult
	urls            []string
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

func NewDisp(w int, wg *sync.WaitGroup, urls []string) *disp {
	return &disp{
		NumberOfWorkers: w,
		wg:              wg,
		worker:          worker{},
		urls:            urls,
		jobChannel:      make(chan Job),
		ResultChannel:   make(chan JobResult),
	}
}

func (d disp) Start() {
	for i := 1; i <= d.NumberOfWorkers; i++ {
		go d.worker.Start3(i, d.wg, d.jobChannel, d.ResultChannel)
	}
	d.wg.Add(d.NumberOfWorkers)

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
