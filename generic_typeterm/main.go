package main

import "fmt"

func main() {
	a := 10
	b := 20
	fmt.Println(Min(a, b))
	type MyInt int
	var myA MyInt = 10
	var myB MyInt = 20
	fmt.Println(Min(myA, myB))
}

// Smimilar to “Embedding and Interfaces”:
// https://golang.org/doc/effective_go.html#embedding
// Just like you can embed a type in a struct, you can also embed an interface in an interface.
// For example, the io.ReadCloser interface is built out of an io.Reader, io.Writer, and an io.Closer:
//
//	type Reader interface {
//		Read(p []byte) (n int, err error)
//	}
//
//	type Writer interface {
//		Write(p []byte) (n int, err error)
//	}
//
//	type Closer interface {
//		Close() error
//	}
//
// // ReadWriter is the interface that combines the Reader, Writer, and Closer interfaces.
//
//	type ReadWriterCloser interface {
//		Reader
//		Writer
//		Closer
//	}
//
// If we try to use Min with a user-defined type whose underlying type is one of the types listed in BuiltInOrdered, we’ll get an error.
// If you want a type term to be valid for any type that has the type term as its underlying type, put a ~ before the type term.
type BuiltInOrdered interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func Min[T BuiltInOrdered](v1, v2 T) T {
	if v1 < v2 {
		return v1
	}
	return v2
}
