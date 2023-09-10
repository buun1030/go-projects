package main

import (
	"fmt"
	"strings"
)

func main() {
	t1 := NewTree(BuiltInOrderable[int])
	t1.Add(10)
	t1.Add(30)
	t1.Add(15)
	fmt.Println(t1.Contains(15))
	fmt.Println(t1.Contains(40))

	t2 := NewTree(OrderPeople)
	t2.Add(Person{"Bob", 30})
	t2.Add(Person{"Maria", 35})
	t2.Add(Person{"Bob", 50})
	fmt.Println(t2.Contains(Person{"Bob", 30}))
	fmt.Println(t2.Contains(Person{"Fred", 25}))

	// after supply a method to NewTree
	t3 := NewTree(Person.Order)
	t3.Add(Person{"Bob", 30})
	t3.Add(Person{"Maria", 35})
	t3.Add(Person{"Bob", 50})
	fmt.Println(t3.Contains(Person{"Bob", 30}))
	fmt.Println(t3.Contains(Person{"Fred", 25}))
}

type BuiltInOrdered interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// The secret is to realize that what our tree needs is a single generic function
// that compares two values and tells us their order:
type OrderableFunc[T any] func(t1, t2 T) int

type Tree[T any] struct {
	f    OrderableFunc[T]
	root *Node[T]
}

type Node[T any] struct {
	val         T
	left, right *Node[T]
}

func NewTree[T any](f OrderableFunc[T]) *Tree[T] {
	return &Tree[T]{
		f: f,
	}
}

// Tree ’s methods are very simple, because they just call Node to do all the real work:
func (t *Tree[T]) Add(v T) {
	t.root = t.root.Add(t.f, v)
}

func (t *Tree[T]) Contains(v T) bool {
	return t.root.Contains(t.f, v)
}

// The Add and Contains methods on Node are very similar to what we’ve seen before (see generic_stack directory).
// The only difference is that the function we are using to order our elements is passed in:
func (n *Node[T]) Add(f OrderableFunc[T], v T) *Node[T] {
	if n == nil {
		return &Node[T]{val: v}
	}
	switch r := f(v, n.val); {
	case r <= -1:
		n.left = n.left.Add(f, v)
	case r >= 1:
		n.right = n.right.Add(f, v)
	}
	return n
}

func (n *Node[T]) Contains(f OrderableFunc[T], v T) bool {
	if n == nil {
		return false
	}
	switch r := f(v, n.val); {
	case r <= -1:
		return n.left.Contains(f, v)
	case r >= 1:
		return n.right.Contains(f, v)
	}
	return true
}

// Now we need a function that matches the OrderedFunc definition.
// By taking advantage of BuiltInOrdered,
// we can write a single function that supports any primitive type:
func BuiltInOrderable[T BuiltInOrdered](t1, t2 T) int {
	if t1 < t2 {
		return -1
	}
	if t1 > t2 {
		return 1
	}
	return 0
}

type Person struct {
	Name string
	Age  int
}

func OrderPeople(p1, p2 Person) int {
	out := strings.Compare(p1.Name, p2.Name)
	if out == 0 {
		out = p1.Age - p2.Age
	}
	return out
}

// Instead of using a function, we can also supply a method to NewTree.
// As we talked about in “Methods Are Functions Too”,
// you can use a method expression to treat a method like a function.
// Let’s do that here. First we write the method:
func (p Person) Order(other Person) int {
	out := strings.Compare(p.Name, other.Name)
	if out == 0 {
		out = p.Age - other.Age
	}
	return out
}
