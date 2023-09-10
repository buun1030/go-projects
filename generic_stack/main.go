package main

import (
	"fmt"
)

// any is equivalent to interface{} (starting with Go 1.18 and later)
// For contain method that have to compare
// Just as interface{} doesn’t say anything, neither does any.
// We can only store values of any type and retrieve them.
// To use ==, we need a different type, in this case, `comparable“
type Stack[T comparable] struct {
	vals []T
}

func (s *Stack[T]) Push(val T) {
	s.vals = append(s.vals, val)
}

// we can’t just return nil, because that’s not a valid value for a value type, like int.
// The easiest way to get a zero value for a generic is to simply declare a variable with var and return it,
// since by definition, var always initializes its variable to the zero value if no other value is assigned.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}
	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}

func (s Stack[T]) Contains(val T) bool {
	for _, v := range s.vals {
		if v == val {
			return true
		}
	}
	return false
}

func main() {
	var intStack Stack[int]
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	v, ok := intStack.Pop()
	fmt.Println(v, ok)
	v, ok = intStack.Pop()
	fmt.Println(v, ok)
	v, ok = intStack.Pop()
	fmt.Println(v, ok)
	v, ok = intStack.Pop()
	fmt.Println(v, ok)
	// The only difference is that when we declare our variable,
	// we include the type that we want to use with our Stack, in this case int.
	// If you try to push a string onto our stack, the compiler will catch it.
	// intStack.Push("nope") // compile error

	var s Stack[int]
	s.Push(10)
	s.Push(20)
	s.Push(30)
	fmt.Println(s.Contains(10))
	fmt.Println(s.Contains(5))
}
