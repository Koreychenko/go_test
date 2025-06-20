package main

import (
	"context"
	"net/http"
	"time"
)

const (
	contextTimeout = 10 * time.Second
	address        = ":9993"
)

func VeryLongOperation(ctx context.Context) ([]byte, error) {
	time.Sleep(2 * time.Second)

	return []byte("Hello"), nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/long-process", func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(request.Context(), contextTimeout)
		defer cancel()

		resultChannel := make(chan []byte, 1)

		go func() {
			data, err := VeryLongOperation(ctx)
			if err != nil {
				return
			}

			select {
			case resultChannel <- data:
			case <-ctx.Done():
			}
		}()

		select {
		case val := <-resultChannel:
			writer.Write(val)
			return
		case <-ctx.Done():
			writer.WriteHeader(http.StatusRequestTimeout)
			return
		}
	})

	http.ListenAndServe(address, mux)
}
