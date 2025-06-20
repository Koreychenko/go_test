# Thread-safe Map Implementation
Create a thread-safe map[string]int using sync.Mutex or sync.RWMutex.

Implement methods:

```
Set(key string, value int)
Get(key string) (int, bool)
Delete(key string)
Keys() []string
```

Make sure the code is safe under concurrent access and validate it using `go run -race`.
