package main

// import "fmt"

// type ImpossiblePrintableInt interface {
// 	int
// 	String() string
// }

// type ImpossibleStruct[T ImpossiblePrintableInt] struct {
// 	val T
// }

// type MyInt int

// func (mi MyInt) String() string {
// 	return fmt.Sprint(mi)
// }

// func main() {
// 	s := ImpossibleStruct[int]{10}
// 	fmt.Println(s)
// 	s2 := ImpossibleStruct[MyInt]{10}
// 	fmt.Println(s2)
// }

// Type Elements Limit Constants
// If you use the Integer interface, the following code will not compile,
// because you cannot assign 1,000 to an 8-bit integer:
// // INVALID!
// func PlusOneThousand[T Integer](in T) T {
//     return in + 1_000
// }
// // VALID
// func PlusOneHundred[T Integer](in T) T {
//     return in + 100
// }

// Functional programming does not work.
// Rather than chaining method calls together,
// you need to either nest function calls or use the much more
// readable approach of invoking the functions one at a time
// and assigning the intermediate values to variables.
// type functionalSlice[T any] []T

// // THIS DOES NOT WORK
// func (fs functionalSlice[T]) Map[E any](f func(T) E) functionalSlice[E] {
//     out := make(functionalSlice[E], len(fs))
//     for i, v := range fs {
//         out[i] = f(v)
//     }
//     return out
// }

// // THIS DOES NOT WORK
// func (fs functionalSlice[T]) Reduce[E any](start E, f func(E, T) E) E {
//     out := start
//     for _, v := range fs {
//         out = f(out, v)
//     }
//     return out
// }

// // which you could use like this (but not work):
// var numStrings = functionalSlice[string]{"1", "2", "3"}
// sum := numStrings.Map(func(s string) int {
//     v, _ := strconv.Atoi(s)
//     return v
// }).Reduce(0, func(acc int, cur int) int {
//     return acc + cur
// })
