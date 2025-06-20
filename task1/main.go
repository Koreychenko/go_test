package main

import (
	"strconv"
	"sync"
)

type SafeMap struct {
	mu      sync.RWMutex
	storage map[string]int
}

func NewSafeMap() SafeMap {
	return SafeMap{
		storage: make(map[string]int),
	}
}

func (m *SafeMap) Set(key string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.storage[key] = value
}

func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.storage[key]

	return value, ok
}

func (m *SafeMap) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.storage, key)
}

func (m *SafeMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, len(m.storage))

	for key := range m.storage {
		keys = append(keys, key)
	}

	return keys
}

func main() {
	safeMap := NewSafeMap()

	wg := sync.WaitGroup{}

	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			in := strconv.Itoa(index)

			safeMap.Set("key"+in, index)
			safeMap.Get("key" + in)
			safeMap.Delete("key" + in)
			safeMap.Keys()
		}(i)
	}

	wg.Wait()
}
