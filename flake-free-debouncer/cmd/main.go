package main

import (
	"fmt"
	"sync"
	"time"
)

type Debouncer struct {
	mu       sync.Mutex
	timer    *time.Timer
	interval time.Duration
}

func NewDebouncer(interval time.Duration) *Debouncer {
	return &Debouncer{interval: interval}
}

// Debounce schedules fn to run after the debounce interval. If called again
// before the interval elapses, the previous scheduled call is cancelled and
// the timer resets — ensuring fn fires only once after the last invocation.
func (d *Debouncer) Debounce(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.interval, fn)
}

func main() {
	callCount := 0
	var mu sync.Mutex

	debouncer := NewDebouncer(100 * time.Millisecond)

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			debouncer.Debounce(func() {
				mu.Lock()
				callCount++
				mu.Unlock()
				fmt.Println("API call fired")
			})
		}()
	}

	wg.Wait()
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	fmt.Printf("Callback fired %d time(s)\n", callCount)
	mu.Unlock()
}
