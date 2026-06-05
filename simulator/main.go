package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
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

func worker(jobs <-chan int, wg *sync.WaitGroup, success *LockedInt, fails *LockedInt) {
	defer wg.Done()
	for range jobs {
		conn, err := net.Dial("tcp", "localhost:8000")
		if err != nil {
			fails.Increment()
			continue
		}

		isError := false

		err = sendPayload(conn, Payload(23))
		if err != nil {
			fails.Increment()
			isError = true
		}

		time.Sleep(50 * time.Millisecond)

		decoder := gob.NewDecoder(conn)
		var serverPayload Payload
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
	flag.Parse()

	successCount := LockedInt{}
	failCount := LockedInt{}

	var wg sync.WaitGroup
	jobs := make(chan int, *total)

	start := time.Now()

	// spawn workers
	for i := 0; i < *maxWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg, &successCount, &failCount)
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
