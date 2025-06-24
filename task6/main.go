package main

import (
	"log/slog"
	"sync"
)

func SafeGo(task func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic recovered", slog.Any("error", err))
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

	SafeGo(nil)

	wg.Wait()
}
