package main

import (
	"fmt"
	"strconv"
)

type MyStruct struct {
	Field string
}

func Map[T any, R any](input []T, mapper func(T) R) []R {
	output := make([]R, 0, len(input))
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

	fmt.Println(intsToStrings)

	extractFields := Map[MyStruct, string](
		[]MyStruct{{"one"}, {"two"}, {"three"}},
		func(i MyStruct) string {
			return i.Field
		},
	)

	fmt.Println(extractFields)
}
