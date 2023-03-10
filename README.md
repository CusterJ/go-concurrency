# go-concurrency
This project implements concurrency in GO through a dispatcher, workers, logger, and benchmark.

There are two types of implemented workers:

    Through a dispatcher that manages workers and has three worker types.
    Worker without a dispatcher.

The dispatcher (or worker) receives a list of links, creates a work channel and a result channel, launches workers that read from the work channel and return the response to the result channel. The logger listens to the result channel and writes the workers' job results to a log file.



![image](https://user-images.githubusercontent.com/104718422/224294044-2ef40552-f097-4a41-92e2-6161bff913bc.png)
