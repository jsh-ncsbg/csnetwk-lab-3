package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/json"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/payload"
)

func worker(jobs <-chan int, delayMs int, sChan chan time.Duration) {
	for range jobs {
		conn, err := net.Dial("tcp", "localhost:8080")

		start := time.Now()
		if err != nil {
			continue
		}

		isError := false
		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		err = json.Write(conn, payload.Payload{Name: "Josh", Number: 1})
		if err != nil {
			isError = true
		}

		var tmp payload.Payload
		err = json.Read(conn, &tmp)
		if err != nil {
			isError = true
		}

		if !isError {
			sChan <- time.Since(start)
		}

		conn.Close()
	}
}

func main() {
	maxWorkers := flag.Int("concurrent", 50, "number of concurrent users")
	total := flag.Int("total", 1000, "total number of users")
	delayMs := flag.Int("delay", 0, "Set time (in milliseconds) before client replies")
	flag.Parse()

	success := 0

	successChan := make(chan time.Duration, *total)

	var wg sync.WaitGroup
	jobs := make(chan int, *total)

	start := time.Now()

	// spawn workers
	for i := 0; i < *maxWorkers; i++ {
		wg.Go(func() {
			worker(jobs, *delayMs, successChan)
		})
	}

	// spawn jobs
	for i := 0; i < *total; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	close(successChan)

	success = len(successChan)
	fail := *total - success

	elapsedTime := time.Since(start)
	fmt.Println("============= BENCHMARK RESULTS =============")
	fmt.Printf("Concurrent Users: %d\n", min(*maxWorkers, *total))
	fmt.Printf("Total Users: %d\n", *total)
	fmt.Printf("Completed in %v\n", elapsedTime.Round(time.Millisecond))
	fmt.Printf("Success: %d/%d\n", success, *total)
	fmt.Printf("Fail: %d/%d\n", fail, *total)
}
