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

func processClient(conn net.Conn) error {
	defer conn.Close()

	// send payload to server
	if err := json.Write(conn, payload.NewPayload("Josh", 67)); err != nil {
		return err
	}

	// receive payload form server
	server := &payload.Payload{}
	if err := json.Read(conn, server); err != nil {
		return err
	}

	return nil
}

func worker(jobs <-chan int, delayMs time.Duration, sChan chan time.Duration) {
	for range jobs {
		conn, err := net.Dial("tcp", "localhost:8080")

		start := time.Now()
		if err != nil {
			continue
		}

		// simulate lag
		time.Sleep(delayMs)

		if err := processClient(conn); err == nil {
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
			worker(jobs, time.Duration(*delayMs)*time.Millisecond, successChan)
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
	fmt.Printf("Time Elapsed: %v\n", elapsedTime.Round(time.Millisecond))
	fmt.Printf("Success: %d (%.2f%%)\n", success, 100 * float64(success)/float64(*total))
	fmt.Printf("Fail: %d (%.2f%%)\n", fail, 100 * float64(fail)/float64(*total))
}
