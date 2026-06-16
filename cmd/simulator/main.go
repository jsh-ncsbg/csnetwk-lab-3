package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jsh-ncsbg/csnetwk-lab-3/protocol"
)

type LockedInt struct {
	mu sync.Mutex
	int
}

func (li *LockedInt) Increment() {
	li.mu.Lock()
	li.int++
	li.mu.Unlock()
}

func worker(jobs <-chan int, delayMs int, wg *sync.WaitGroup, success *LockedInt, fails *LockedInt) {
	defer wg.Done()
	for range jobs {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			fails.Increment()
			continue
		}

		isError := false

		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		err = protocol.SendPayload(conn, protocol.Payload(23))
		if err != nil {
			fails.Increment()
			isError = true
		}

		decoder := gob.NewDecoder(conn)
		var serverPayload protocol.Payload
		if err := decoder.Decode(&serverPayload); err != nil {
			fails.Increment()
			isError = true
		}

		if !isError {
			success.Increment()
		}

		conn.Close()
	}
}

func main() {
	maxWorkers := flag.Int("concurrent", 50, "number of concurrent users")
	total := flag.Int("total", 1000, "total number of users")
	delayMs := flag.Int("delay", 0, "Set time (in milliseconds) before client replies")
	flag.Parse()

	successCount := LockedInt{}
	failCount := LockedInt{}

	var wg sync.WaitGroup
	jobs := make(chan int, *total)

	start := time.Now()

	// spawn workers
	for i := 0; i < *maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, *delayMs, &wg, &successCount, &failCount)
	}

	// spawn jobs
	for i := 0; i < *total; i++ {
		jobs <- i
	}
	close(jobs)

	wg.Wait()

	elapsedTime := time.Since(start)
	fmt.Println("============= BENCHMARK RESULTS =============")
	fmt.Printf("Concurrent Users: %d\n", min(*maxWorkers, *total))
	fmt.Printf("Total Users: %d\n", *total)
	fmt.Printf("Completed in %.2f seconds\n", elapsedTime.Seconds())
	fmt.Printf("Success: %d/%d\n", successCount.int, *total)
	fmt.Printf("Fail: %d/%d\n", failCount.int, *total)
}
