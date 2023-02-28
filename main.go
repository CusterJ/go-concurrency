package main

import (
	"concurrency/links"
	"concurrency/logger"
	"concurrency/workers"
	"log"
	"sync"
)

func main() {
	lks := links.GenerateLinks(200)
	log.Println(lks)

	wg := &sync.WaitGroup{}
	disp := workers.NewDisp(10, wg, lks)

	lg := logger.NewLog(disp.ResultChannel)
	lg.CreateLogFIle()

	done := make(chan bool)

	lg.LogStart(done)
	disp.Start()

	wg.Wait()
	<-done

	lg.CloseLogFIle()
}
