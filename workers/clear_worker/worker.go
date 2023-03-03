package clear_worker

import (
	"net/http"
	"sync"
	"time"
)

type worker struct {
	id            int
	wg            *sync.WaitGroup
	jobChannel    chan job
	ResultChannel chan JobResult
	urls          []string
}

type job struct {
	id   int
	site string
}

type JobResult struct {
	WorkerID int
	JobID    int
	Site     string
	Time     float64
	Response Response
}

type Response struct {
	Status int
	Err    string
}

// Takes a slice of links as work and returns a worker instance.
func NewWorker(urls []string) *worker {
	return &worker{
		urls:          urls,
		ResultChannel: make(chan JobResult, len(urls)),
	}
}

// Starts the specified number of workers.
// Creates a sync.WaitGroup{}, work channel, fills the work channel.
// Waits for work to be done and closes the work channel and the job result channel.
func (w *worker) Start(workers int) {
	w.wg = &sync.WaitGroup{}
	w.jobChannel = make(chan job, len(w.urls))

	for i := 1; i <= workers; i++ {
		go worker.new(*w, i)
	}
	w.wg.Add(workers)

	for i, v := range w.urls {
		w.jobChannel <- job{
			id:   i,
			site: v,
		}
	}
	close(w.jobChannel)
	w.wg.Wait()
	close(w.ResultChannel)
}

func (w worker) new(id int) {
	defer w.wg.Done()
	w.id = id
	response := Response{}
	// log.Printf("Worker %d started", w.id)

	for j := range w.jobChannel {
		start := time.Now()

		status, err := goToWebsite(j.site)
		if err != nil {
			response = Response{
				Err: err.Error(),
			}
		} else {
			response = Response{
				Status: status,
			}
		}

		t := time.Since(start)

		w.ResultChannel <- JobResult{
			WorkerID: w.id,
			JobID:    j.id,
			Site:     j.site,
			Time:     t.Seconds(),
			Response: response,
		}
		// s := fmt.Sprintf("worker id: %3d job_id: %3d %-30s| %6.2fs | %4d %s", w.id, j.id, j.site, t.Seconds(), response.Status, response.Err)
		// log.Println(s)
	}
	// log.Println("Stop worker: ", w.id)
}

func goToWebsite(site string) (status int, err error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	res, err := client.Get(site)
	if err != nil {
		return status, err
	}
	status = res.StatusCode
	return status, nil
}
