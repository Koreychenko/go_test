package main

import (
	"fmt"
	"log/slog"
	"sync"
)

func SafeGo(task func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if recoverError, ok := err.(error); ok {
					slog.Error(fmt.Sprintf("panic: %s", recoverError.Error()))
				} else {
					slog.Error(fmt.Sprintf("panic: %s", err))
				}
			}
		}()

		task()
	}()
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	SafeGo(func() {
		defer wg.Done()
		panic("PAAANIC")
	})

	wg.Add(1)
	SafeGo(func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			print(i)
		}
	})

	wg.Wait()
}
