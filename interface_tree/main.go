package main

import "strings"

func main() {
	var it1 *Tree
	it1 = it1.Insert(OrderableInt(5))
	it1 = it1.Insert(OrderableInt(3))
	// etc...

	var it2 *Tree
	// it2 = it2.Insert(OrderableInt(5))
	it2 = it2.Insert(OrderableString("nope"))
}

type Orderable interface {
	// Order returns:
	// a value < 0 when the Orderable is less than the supplied value,
	// a value > 0 when the Orderable is greater than the supplied value,
	// and 0 when the two values are equal.
	Order(interface{}) int
}

type Tree struct {
	val         Orderable
	left, right *Tree
}

func (t *Tree) Insert(val Orderable) *Tree {
	if t == nil {
		return &Tree{val: val}
	}

	switch comp := val.Order(t.val); {
	case comp < 0:
		t.left = t.left.Insert(val)
	case comp > 0:
		t.right = t.right.Insert(val)
	}
	return t
}

type OrderableInt int

func (oi OrderableInt) Order(val interface{}) int {
	return int(oi - val.(OrderableInt))
}

type OrderableString string

func (os OrderableString) Order(val interface{}) int {
	return strings.Compare(string(os), val.(string))
}
