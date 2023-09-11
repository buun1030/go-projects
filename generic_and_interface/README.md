Exercises
Now that you’ve seen how generics work, apply them to the solve the following problems. Solutions are available in the exercise_solutions directory in the [chapter 8 repository](https://github.com/learning-go-book-2e/ch08).

Write a generic function that doubles the value of any integer or float that’s passed in to it. Define any needed generic interfaces.

Define a generic interface called Printable that matches a type that implements fmt.Stringer and has an underlying type of int or float64. Define types that meet this interface. Write a function that takes in a Printable and prints its value to the screen using fmt.Println.

Write a generic singly linked list data type. Each element can hold a comparable value and has a pointer to the next element in the list. The methods to implement are:

// adds a new element to the end of the linked list
Add(T)
// adds an element at the specified position in the linked list
Insert(T, int)
// returns the position of the supplied value, -1 if it's not present
Index (T) int