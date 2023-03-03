package main

import (
	"concurrency/links"
	"concurrency/logger"
	worker "concurrency/workers/clear_worker"
	"log"
)

func main() {
	links := links.GenerateLinks(100)

	// WORKER WITH DISPATCHER
	// disp := workers.NewDisp(links)
	// logger := logger.NewLog(disp.ResultChannel)
	// disp.Start(10)

	// WORKER WITHOUT DISPATCHER
	worker := worker.NewWorker(links)
	logger := logger.NewLogger(worker.ResultChannel)

	logger.LogStart()
	worker.Start(10)

	<-logger.DoneCh
	log.Println("main func end")
}
