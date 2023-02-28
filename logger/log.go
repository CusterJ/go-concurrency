package logger

import (
	"concurrency/workers"
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	file      *os.File
	ch        chan workers.JobResult
	timeStart time.Time
	DoneCh    chan bool
}

func NewLog(ch chan workers.JobResult) *Log {
	return &Log{
		ch:     ch,
		DoneCh: make(chan bool),
	}
}

func (l *Log) CreateLogFIle() {
	file, err := os.Create("Log.txt")
	if err != nil {
		log.Println("func LogFile -> can't create Log.txt: ", err)
		return
	}
	l.file = file
	log.Println("Create file: ", l.file.Name())
}

func (l *Log) LogStart() {
	l.CreateLogFIle()

	l.timeStart = time.Now()
	workers := make(map[int]int)
	responses := make(map[string]int)
	var jobs int

	go func() {
		for jr := range l.ch {
			s := fmt.Sprintf("Worker id: %3d Job_id: %3d %-30s| %6.2fs | %s\n", jr.WorkerID, jr.JobID, jr.Site, jr.Time, jr.Response)

			_, err := l.file.WriteString(s)
			if err != nil {
				log.Println("LogFromChannel => error Write string: ", err)
			}
			//count TOTALS
			jobs++
			workers[jr.WorkerID]++
			responses[jr.Response]++
		}

		//Write total to Log File
		t := fmt.Sprintf("=== TOTALS ===\nJobs Done %d\n", jobs)
		l.file.WriteString(t)

		l.file.WriteString("WORKERS: JOBS\n")
		for i, v := range workers {
			str := fmt.Sprintf("Worker_%-3d: %d\n", i, v)
			l.file.WriteString(str)
		}

		l.file.WriteString("RESPONSES: STATUSES\n")
		for i, v := range responses {
			str := fmt.Sprintf("%3d: %s\n", v, i)
			l.file.WriteString(str)
		}
		log.Println("LogFromChannel => Result channel closed => stop logging")

		//Add final string and close log file
		l.StopLog()
		l.DoneCh <- true
	}()
}

func (l *Log) StopLog() {
	defer l.file.Close()

	t := time.Since(l.timeStart).Seconds()
	s := fmt.Sprintf("TOTAL TIME log file: %.2f", t)
	_, err := l.file.WriteString(s)
	if err != nil {
		log.Println("LogFromChannel => error Write string: ", err)
	}
	log.Println("Close file: ", l.file.Name())
}
