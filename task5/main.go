package main

import (
	"fmt"
	"strconv"
)

type Str struct {
	Field string
}

func Map[T any, R any](input []T, mapper func(T) R) []R {
	var output []R
	for _, val := range input {
		output = append(output, mapper(val))
	}

	return output
}

func main() {
	intsToStrings := Map[int, string](
		[]int{1, 2, 3, 4},
		func(i int) string {
			return strconv.Itoa(i)
		},
	)

	fmt.Printf("%v", intsToStrings)

	extractFields := Map[Str, string](
		[]Str{{"one"}, {"two"}, {"three"}},
		func(i Str) string {
			return i.Field
		},
	)

	fmt.Printf("%v", extractFields)
}
