package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Statistics struct {
	mu          sync.Mutex
	StatusCodes map[string]int
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Starting batch")

	URLs := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
		"https://www.reddit.com",
		"https://www.twitter.com",
		"https://www.youtube.com",
		"https://www.upscolled.com",
		"https://www.apple.com",
		"https://www.orange.com",
	}

	jobs := make(chan Job, 5)

	statistics := Statistics{
		StatusCodes: make(map[string]int),
	}

	numOfWorkers := 3
	var wg sync.WaitGroup

	for i := range numOfWorkers {
		wg.Add(1)
		go worker(ctx, i, jobs, &statistics, &wg)
	}

	for _, url := range URLs {
		jobs <- Job{URL: url}
	}

	close(jobs)

	wg.Wait()

	fmt.Println("\n--- Final Scraping Statistics ---")
	for url, statusCode := range statistics.StatusCodes {
		fmt.Printf("%s -> HTTP %d\n", url, statusCode)
	}
}

func worker(ctx context.Context, workerID int, jobs <-chan Job, statistics *Statistics, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Starting worker %d\n", workerID)

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Printf("Finished worker %d\n", workerID)
				return
			}
			process(job, statistics)
		case <-ctx.Done():
			fmt.Printf("Stopping worker %d\n", workerID)
			return
		}
	}
}

func process(job Job, statistics *Statistics) {
	url := job.URL
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)

	statusCode := 0
	if err != nil {
		fmt.Printf("Error fetching %s: %s\n", url, err)
	} else {
		statusCode = resp.StatusCode
		resp.Body.Close()
	}

	statistics.mu.Lock()
	statistics.StatusCodes[url] = statusCode
	statistics.mu.Unlock()
}
