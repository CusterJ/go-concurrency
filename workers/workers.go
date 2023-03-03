package workers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type worker struct {
	id int
}

func (w worker) Start(id int, wg *sync.WaitGroup, jobChannel <-chan Job, rc chan<- string) {
	w.id = id
	// log.Printf("Worker %d started", w.id)

	go func() {
		defer wg.Done()
		client := http.Client{
			Timeout: time.Second * 10,
		}
		for {
			j, ok := <-jobChannel
			if ok {
				start := time.Now()

				var result string

				res, err := client.Get(j.Site)
				if err != nil {
					result = "error"
				} else {
					result = res.Status
				}

				t := time.Since(start)
				s := fmt.Sprintf("worker id: %3d job_id: %3d %-30s| %6.2fs | %s", w.id, j.ID, j.Site, t.Seconds(), result)
				// log.Println(s)
				rc <- s
			} else {
				// log.Println("Stop worker: ", w.id)
				return
			}
		}
	}()
}

func (w worker) Start2(id int, wg *sync.WaitGroup, jobChannel <-chan Job, rc chan<- string) {
	defer wg.Done()
	w.id = id
	// log.Printf("Worker %d started", w.id)

	for {
		j, ok := <-jobChannel
		if ok {
			start := time.Now()
			res := w.doJob(j.Site)
			t := time.Since(start)
			s := fmt.Sprintf("worker id: %3d job_id: %3d %-30s| %6.2fs | %s", w.id, j.ID, j.Site, t.Seconds(), res)
			// log.Println(s)

			rc <- s
		} else {
			// log.Println("Stop worker: ", w.id)
			return
		}
	}
}

func (w *worker) doJob(url string) string {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "error"
	}
	defer res.Body.Close()

	return res.Status
}

func (w worker) Start3(id int, wg *sync.WaitGroup, jobChannel <-chan Job, rc chan<- JobResult) {
	defer wg.Done()
	w.id = id
	// log.Printf("Worker %d started", w.id)
	client := http.Client{
		Timeout: time.Second * 10,
	}

	for j := range jobChannel {
		start := time.Now()

		var result string

		res, err := client.Get(j.Site)
		if err != nil {
			result = err.Error()
		} else {
			result = strconv.Itoa(res.StatusCode)
		}

		t := time.Since(start)
		// s := fmt.Sprintf("worker id: %3d job_id: %3d %-30s| %6.2fs | %s", w.id, j.ID, j.Site, t.Seconds(), result)
		// log.Println(s)

		rc <- JobResult{
			WorkerID: w.id,
			JobID:    j.ID,
			Site:     j.Site,
			Time:     t.Seconds(),
			Response: result,
		}
	}
	// log.Println("Stop worker: ", w.id)
}
