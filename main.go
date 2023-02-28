package main

import (
	"concurrency/links"
	"concurrency/logger"
	"concurrency/workers"
	"log"
)

func main() {
	lks := links.GenerateLinks(20)
	log.Println(lks)

	disp := workers.NewDisp(lks)

	lg := logger.NewLog(disp.ResultChannel)

	lg.LogStart()

	disp.Start(10)
	<-lg.DoneCh
	log.Println("main func")
}
