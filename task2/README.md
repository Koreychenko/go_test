# Worker Pool with Goroutines
Implement a worker pool where:

Workers consume tasks from a channel,

Process the tasks concurrently,

Send results to another channel.

Support graceful shutdown using context.Context and sync.WaitGroup.
