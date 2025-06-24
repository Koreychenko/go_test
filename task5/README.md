# Generics: Map Function
Implement a generic Map[T, R] function:

```go
func Map[T any, R any](input []T, mapper func(T) R) []R
```

Use it with:

Primitive types (`[]int → []string`),

Custom struct types.

Generics are essential in modern Go (1.18+).
