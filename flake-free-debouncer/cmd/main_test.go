package main

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

func TestDebouncer_FiresOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		callCount := 0
		var mu sync.Mutex

		debouncer := NewDebouncer(100 * time.Millisecond)

		// Spawn multiple goroutines all triggering the debouncer concurrently.
		var wg sync.WaitGroup
		for range 10 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				debouncer.Debounce(func() {
					mu.Lock()
					callCount++
					mu.Unlock()
				})
			}()
		}

		// Wait for all goroutines to finish scheduling, then advance fake time
		// past the debounce interval so the timer fires.
		wg.Wait()
		synctest.Wait()
		time.Sleep(200 * time.Millisecond)
		synctest.Wait()

		mu.Lock()
		got := callCount
		mu.Unlock()

		if got != 1 {
			t.Errorf("expected callback to fire 1 time, got %d", got)
		}
	})
}

func TestDebouncer_ResetsOnSubsequentCall(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		callCount := 0
		var mu sync.Mutex

		debouncer := NewDebouncer(100 * time.Second)

		// First burst — should be debounced to a single fire.
		var wg sync.WaitGroup
		for range 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				debouncer.Debounce(func() {
					mu.Lock()
					callCount++
					mu.Unlock()
				})
			}()
		}
		wg.Wait()
		synctest.Wait()
		time.Sleep(200 * time.Second)
		synctest.Wait()

		// Second burst — should produce exactly one more fire.
		for range 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				debouncer.Debounce(func() {
					mu.Lock()
					callCount++
					mu.Unlock()
				})
			}()
		}
		wg.Wait()
		synctest.Wait()
		time.Sleep(200 * time.Second)
		synctest.Wait()

		mu.Lock()
		got := callCount
		mu.Unlock()

		if got != 2 {
			t.Errorf("expected callback to fire 2 times (once per burst), got %d", got)
		}
	})
}
