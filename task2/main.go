package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	maxWorkers = 5
	maxTasks   = 1000
	timeout    = 50 * time.Second
)

func worker(ctx context.Context, inCh <-chan int, outCh chan<- int) {
	for {
		select {
		case value, ok := <-inCh:
			if !ok {
				return
			}

			time.Sleep(1 * time.Second)
			outCh <- value
		case <-ctx.Done():
			return
		}
	}
}

func feeder(ctx context.Context, maxTasks int, inCh chan<- int) {
	defer close(inCh)
	for i := 0; i < maxTasks; i++ {
		select {
		case inCh <- i:
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	inCh := make(chan int, maxWorkers)
	outCh := make(chan int, maxWorkers)
	var wg sync.WaitGroup

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		for i := 0; i < maxWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				worker(ctx, inCh, outCh)
			}()
		}
		wg.Wait()

		close(outCh)
	}()

	go feeder(ctx, maxTasks, inCh)

	for val := range outCh {
		fmt.Printf("%d ", val)
	}
}
