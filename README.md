# go-mistakes

## Unintended variable shadowing

Declares outer block variable and then redeclare the same name variable in an inner block,
after prcess the inner block, the outer block variable will stay the same.\
\
**Solution:** To assign value to the outer block variable in the inner block, just use `=` not `:=`

## Unnecessary nested code

Because it was difficult to distinguish the expected execution flow because of the nested if/else statements. Conversely, if only one nested if/else statment, it requires scanning down one column to see the expected execution flow and down the second column to see how the edge cases are handled.\
\
![This is an alt text.](https://drek4537l1klr.cloudfront.net/harsanyi/Figures/CH02_F01_Harsanyi.png "To understand the expected execution flow, we just have to scan the happy path column.")
\
\
**Solution:** Striving to reduce the number of nested blocks, aligning the happy path on the left, and returning as early as possible are concrete means to improve our code’s readability.

## Interface on the producer side

*abstractions should be discovered, not created.* This means that it’s not up to the producer to force a given abstraction for all the clients. Instead, it’s up to the client to decide whether it needs some form of abstraction and then determine the best abstraction level for its needs.
\
**Solution:** An interface should live on the consumer side in most cases. However, in particular contexts (for example, when we know—not foresee—that an abstraction will be helpful for consumers), we may want to have it on the producer side. If we do, we should strive to keep it as minimal as possible, increasing its reusability potential and making it more easily composable.
\
**Note:** The interface in Go is satisfied implicitly while some is an explicit implementation which have to declare that a particular class or type explicitly implements a specific interface and cannot do like solution above.

## Not knowing which type of receiver to use

### A receiver must be a pointer

* If the method needs to mutate the receiver. This rule is also valid if the receiver is a slice and a method needs to append elements.
* If the method receiver contains a field that cannot be copied: for example, a type part of the `sync` package.

### A receiver should be a pointer

* If the receiver is a large object. Using a pointer can make the call more efficient, as doing so prevents making an extensive copy. When in doubt about how large is large, benchmarking can be the solution; it’s pretty much impossible to state a specific size, because it depends on many factors.

### A receiver must be a value

* If we have to enforce a receiver’s immutability.
* If the receiver is a map, function, or channel. Otherwise, a compilation error occurs.

### A receiver should be a value

* If the receiver is a slice that doesn’t have to be mutated.
* If the receiver is a small array or struct that is naturally a value type without mutable fields, such as `time.Time`.
* If the receiver is a basic type such as `int`, `float64`, or `string`.

## The standard library

## Not using testing utility packages

## Not exploring all the go testing features
