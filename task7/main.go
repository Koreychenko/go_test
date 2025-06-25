package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	Data string
}

func ChannelListener(taskCh <-chan Task, shutdown <-chan struct{}, heartbeat <-chan time.Time) {
	for {
		select {
		case task := <-taskCh:
			fmt.Println(task.Data)
		case <-shutdown:
			fmt.Println("exit on shutdown")
			return
		case <-heartbeat:
			fmt.Println("exit on timeout")
			return
		}
	}
}

func main() {
	taskCh := make(chan Task, 1)
	shutdownCh := make(chan struct{}, 1)
	heartbeatCh := time.After(5 * time.Second)

	automaticShutdownTimer := time.NewTimer(3 * time.Second)
	go func() {
		select {
		case <-automaticShutdownTimer.C:
			shutdownCh <- struct{}{}
		}
	}()

	var wg sync.WaitGroup

	listenerCloseSignal := make(chan struct{}, 1)

	wg.Add(1)

	go func() {
		defer wg.Done()
		ChannelListener(taskCh, shutdownCh, heartbeatCh)
		listenerCloseSignal <- struct{}{}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(taskCh)
		for i := 1; i <= 10; i++ {
			select {
			case <-listenerCloseSignal:
				return
			default:
				taskCh <- Task{Data: "test" + strconv.Itoa(i)}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait()
}
