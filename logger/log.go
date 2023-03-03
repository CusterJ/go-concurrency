package logger

import (
	"concurrency/workers/clear_worker"
	"fmt"
	"log"
	"os"
	"time"
)

type Log struct {
	file      *os.File
	ch        chan clear_worker.JobResult
	timeStart time.Time
	DoneCh    chan bool
}

// type Log struct {
// 	file      *os.File
// 	ch        chan workers.JobResult
// 	timeStart time.Time
// 	DoneCh    chan bool
// }

// func NewLog(ch chan workers.JobResult) *Log {
// 	return &Log{
// 		ch:     ch,
// 		DoneCh: make(chan bool),
// 	}
// }

// Takes a work channel as an argument and returns an instance of the logger.
func NewLogger(ch chan clear_worker.JobResult) *Log {
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

// Starts the logger and returns a channel to wait for the logger to complete and write the log file.
func (l *Log) LogStart() {
	l.CreateLogFIle()

	l.timeStart = time.Now()
	workers := make(map[int]int)
	statusCode := make(map[int]int)
	errors := make(map[string]int)
	var jobs int

	go func() {
		for jr := range l.ch {
			s := fmt.Sprintf("Worker id: %3d Job_id: %3d %-30s| %6.2fs | %5d %s\n", jr.WorkerID, jr.JobID, jr.Site, jr.Time, jr.Response.Status, jr.Response.Err)

			_, err := l.file.WriteString(s)
			if err != nil {
				log.Println("LogFromChannel => error Write string: ", err)
			}
			//count TOTALS
			jobs++
			workers[jr.WorkerID]++
			if jr.Response.Status != 0 {
				statusCode[jr.Response.Status]++
			}
			if jr.Response.Err != "" {
				errors[jr.Response.Err]++
			}
		}

		//Write total to Log File
		t := fmt.Sprintf("\n=== TOTALS ===\nJobs Done %d\n", jobs)
		l.file.WriteString(t)

		l.file.WriteString("\nWORKERS   : JOBS\n")
		for i, v := range workers {
			str := fmt.Sprintf("Worker_%-3d: %d\n", i, v)
			l.file.WriteString(str)
		}

		l.file.WriteString("\nRESPONSES: STATUSES\n")
		for i, v := range statusCode {
			str := fmt.Sprintf("%-8d : %d\n", v, i)
			l.file.WriteString(str)
		}

		l.file.WriteString("\nRESPONSES: ERRORS\n")
		for i, v := range errors {
			str := fmt.Sprintf("%-8d : %s\n", v, i)
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
	s := fmt.Sprintf("\nTOTAL TIME log file: %.2f", t)
	_, err := l.file.WriteString(s)
	if err != nil {
		log.Println("LogFromChannel => error Write string: ", err)
	}
	log.Println("Close file: ", l.file.Name())
}
