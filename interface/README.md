A Quick Lesson on Interfaces
While Go’s concurrency model (which we cover in Chapter 10) gets all of the publicity, the real star of Go’s design is its implicit interfaces, the only abstract type in Go. Let’s see what makes them so great.

We’ll start by taking a quick look at how to declare interfaces. At their core, interfaces are simple. Like other user-defined types, you use the type keyword.

Here’s the definition of the Stringer interface in the fmt package:

type Stringer interface {
    String() string
}
In an interface declaration, an interface literal appears after the name of the interface type. It lists the methods that must be implemented by a concrete type to meet the interface. The methods defined by an interface are called the method set of the interface.

Like other types, interfaces can be declared in any block.

Interfaces are usually named with “er” endings. We’ve already seen fmt.Stringer, but there are many more, including io.Reader, io.Closer, io.ReadCloser, json.Marshaler, and http.Handler.

Interfaces Are Type-Safe Duck Typing
So far, nothing that’s been said is much different from interfaces in other languages. What makes Go’s interfaces special is that they are implemented implicitly. A concrete type does not declare that it implements an interface. If the method set for a concrete type contains all of the methods in the method set for an interface, the concrete type implements the interface. This means that the concrete type can be assigned to a variable or field declared to be of the type of the interface.

This implicit behavior makes interfaces the most interesting thing about types in Go, because they enable both type-safety and decoupling, bridging the functionality in both static and dynamic languages.

To understand why, let’s talk about why languages have interfaces. Earlier we mentioned that Design Patterns taught developers to favor composition over inheritance. Another piece of advice from the book is “Program to an interface, not an implementation.” Doing so allows you to depend on behavior, not on implementation, allowing you to swap implementations as needed. This allows your code to evolve over time, as requirements inevitably change.

Dynamically typed languages like Python, Ruby, and JavaScript don’t have interfaces. Instead, those developers use “duck typing,” which is based on the expression “If it walks like a duck and quacks like a duck, it’s a duck.” The concept is that you can pass an instance of a type as a parameter to a function as long as the function can find a method to invoke that it expects:

class Logic:
def process(self, data):
    # business logic

def program(logic):
    # get data from somewhere
    logic.process(data)

logicToUse = Logic()
program(logicToUse)
Duck typing might sound weird at first, but it’s been used to build large and successful systems. If you program in a statically typed language, this sounds like utter chaos. Without an explicit type being specified, it’s hard to know exactly what functionality should be expected. As new developers move on to a project or the existing developers forget what the code is doing, they have to trace through the code to figure out what the actual dependencies are.

Java developers use a different pattern. They define an interface, create an implementation of the interface, but only refer to the interface in the client code:

public interface Logic {
    String process(String data);
}

public class LogicImpl implements Logic {
    public String process(String data) {
        // business logic
    }
}

public class Client {
    private final Logic logic;
    // this type is the interface, not the implementation

    public Client(Logic logic) {
        this.logic = logic;
    }

    public void program() {
        // get data from somewhere
        this.logic.process(data);
    }
}

public static void main(String[] args) {
    Logic logic = new LogicImpl();
    Client client = new Client(logic);
    client.program();
}
Dynamic language developers look at the explicit interfaces in Java and don’t see how you can possibly refactor your code over time when you have explicit dependencies. Switching to a new implementation from a different provider means rewriting your code to depend on a new interface.

Go’s developers decided that both groups are right. If your application is going to grow and change over time, you need flexibility to change implementation. However, in order for people to understand what your code is doing (as new people work on the same code over time), you also need to specify what the code depends on. That’s where implicit interfaces come in. Go code is a blend of the previous two styles:

type LogicProvider struct {}

func (lp LogicProvider) Process(data string) string {
    // business logic
}

type Logic interface {
    Process(data string) string
}

type Client struct{
    L Logic
}

func(c Client) Program() {
    // get data from somewhere
    c.L.Process(data)
}

main() {
    c := Client{
        L: LogicProvider{},
    }
    c.Program()
}
In the Go code, there is an interface, but only the caller (Client) knows about it; there is nothing declared on LogicProvider to indicate that it meets the interface. This is sufficient to both allow a new logic provider in the future and provide executable documentation to ensure that any type passed into the client will match the client’s need.

TIP
Interfaces specify what callers need. The client code defines the interface to specify what functionality it requires.

This doesn’t mean that interfaces can’t be shared. We’ve already seen several interfaces in the standard library that are used for input and output. Having a standard interface is powerful; if you write your code to work with io.Reader and io.Writer, it will function correctly whether it is writing to a file on local disk or a value in memory.

Furthermore, using standard interfaces encourages the decorator pattern. It is common in Go to write factory functions that take in an instance of an interface and return another type that implements the same interface. For example, say you have a function with the following definition:

func process(r io.Reader) error
You can process data from a file with the following code:

r, err := os.Open(fileName)
if err != nil {
    return err
}
defer r.Close()
return process(r)
The os.File instance returned by os.Open meets the io.Reader interface and can be used in any code that reads in data. If the file is gzip-compressed, you can wrap the io.Reader in another io.Reader:

r, err := os.Open(fileName)
if err != nil {
    return err
}
defer r.Close()
gz, err = gzip.NewReader(r)
if err != nil {
    return err
}
defer gz.Close()
return process(gz)
Now the exact same code that was reading from an uncompressed file is reading from a compressed file instead.

TIP
If there’s an interface in the standard library that describes what your code needs, use it!

It’s perfectly fine for a type that meets an interface to specify additional methods that aren’t part of the interface. One set of client code may not care about those methods, but others do. For example, the io.File type also meets the io.Writer interface. If your code only cares about reading from a file, use the io.Reader interface to refer to the file instance and ignore the other methods.

Embedding and Interfaces
Just like you can embed a type in a struct, you can also embed an interface in an interface. For example, the io.ReadCloser interface is built out of an io.Reader and an io.Closer:

type Reader interface {
        Read(p []byte) (n int, err error)
}

type Closer interface {
        Close() error
}

type ReadCloser interface {
        Reader
        Closer
}
NOTE
Just like you can embed a concrete type in a struct, you can also embed an interface in a struct. We’ll see a use for this in “Stubs in Go”.

Accept Interfaces, Return Structs
You’ll often hear experienced Go developers say that your code should “Accept interfaces, return structs.” What this means is that the business logic invoked by your functions should be invoked via interfaces, but the output of your functions should be a concrete type. We’ve already covered why functions should accept interfaces: they make your code more flexible and explicitly declare exactly what functionality is being used.

If you create an API that returns interfaces, you are losing one of the main advantages of implicit interfaces: decoupling. You want to limit the third-party interfaces that your client code depends on because your code is now permanently dependent on the module that contains those interfaces, as well as any dependencies of that module, and so on. (We talk about modules and dependencies in Chapter 9.) This limits future flexibility. To avoid the coupling, you’d have to write another interface and do a type conversion from one to the other. While depending on concrete instances can lead to dependencies, using a dependency injection layer in your application limits the effect. We’ll talk more about dependency injection in “Implicit Interfaces Make Dependency Injection Easier”.

Another reason to avoid returning interfaces is versioning. If a concrete type is returned, new methods and fields can be added without breaking existing code. The same is not true for an interface. Adding a new method to an interface means that you need to update all existing implementations of the interface, or your code breaks. If you make a backward-breaking change to an API, you should increment your major version number.

Rather than writing a single factory function that returns different instances behind an interface based on input parameters, try to write separate factory functions for each concrete type. In some situations (such as a parser that can return one or more different kinds of tokens), it’s unavoidable and you have no choice but to return an interface.

Errors are an exception to this rule. As we’ll see in Chapter 8, Go functions and methods declare a return parameter of the error interface type. In the case of error, it’s quite likely that different implementation of the interface could be returned, so you need to use an interface to handle all possible options, as interfaces are the only abstract type in Go.

There is one potential drawback to this pattern. As we discussed in “Reducing the Garbage Collector’s Workload”, reducing heap allocations improves performance by reducing the amount of work for the garbage collector. Returning a struct avoids a heap allocation, which is good. However, when invoking a function with parameters of interface types, a heap allocation occurs for each of the interface parameters. Figuring out the trade-off between better abstraction and better performance is something that should be done over the life of your program. Write your code so that it is readable and maintainable. If you find that your program is too slow and you have profiled it and you have determined that the performance problems are due to a heap allocation caused by an interface parameter, then you should rewrite the function to use a concrete type parameter. If multiple implementations of an interface are passed into the function, this will mean creating multiple functions with repeated logic.

Interfaces and nil
When discussing pointers in Chapter 6, we also talked about nil, the zero value for pointer types. We also use nil to represent the zero value for an interface instance, but it’s not as simple as it is for concrete types.

In order for an interface to be considered nil both the type and the value must be nil. The following code prints out true on the first two lines and false on the last:

var s *string
fmt.Println(s == nil) // prints true
var i interface{}
fmt.Println(i == nil) // prints true
i = s
fmt.Println(i == nil) // prints false
You can run it for yourself on The Go Playground.

In the Go runtime, interfaces are implemented as a pair of pointers, one to the underlying type and one to the underlying value. As long as the type is non-nil, the interface is non-nil. (Since you cannot have a variable without a type, if the value pointer is non-nil, the type pointer is always non-nil.)

What nil indicates for an interface is whether or not you can invoke methods on it. As we covered earlier, you can invoke methods on nil concrete instances, so it makes sense that you can invoke methods on an interface variable that was assigned a nil concrete instance. If an interface is nil, invoking any methods on it triggers a panic (which we’ll discuss in “panic and recover”). If an interface is non-nil, you can invoke methods on it. (But note that if the value is nil and the methods of the assigned type don’t properly handle nil, you could still trigger a panic.)

Since an interface instance with a non-nil type is not equal to nil, it is not straightforward to tell whether or not the value associated with the interface is nil when the type is non-nil. You must use reflection (which we’ll discuss in “Use Reflection to Check If an Interface’s Value Is nil”) to find out.

The Empty Interface Says Nothing
Sometimes in a statically typed language, you need a way to say that a variable could store a value of any type. Go uses interface{} to represent this:

var i interface{}
i = 20
i = "hello"
i = struct {
    FirstName string
    LastName string
} {"Fred", "Fredson"}
You should note that interface{} isn’t special case syntax. An empty interface type simply states that the variable can store any value whose type implements zero or more methods. This just happens to match every type in Go. Because an empty interface doesn’t tell you anything about the value it represents, there isn’t a lot you can do with it. One common use of the empty interface is as a placeholder for data of uncertain schema that’s read from an external source, like a JSON file:

// one set of braces for the interface{} type,
// the other to instantiate an instance of the map
data := map[string]interface{}{}
contents, err := ioutil.ReadFile("testdata/sample.json")
if err != nil {
    return err
}
defer contents.Close()
json.Unmarshal(contents, &data)
// the contents are now in the data map
Another use of interface{} is as a way to store a value in a user-created data structure. This is due to Go’s current lack of user-defined generics. If you need a data structure beyond a slice, array, or map, and you don’t want it to only work with a single type, you need to use a field of type interface{} to hold its value. You can try the following code on The Go Playground:

type LinkedList struct {
    Value interface{}
    Next    *LinkedList
}

func (ll *LinkedList) Insert(pos int, val interface{}) *LinkedList {
    if ll == nil || pos == 0 {
        return &LinkedList{
            Value: val,
            Next:    ll,
        }
    }
    ll.Next = ll.Next.Insert(pos-1, val)
    return ll
}
WARNING
This is not an efficient implementation of insert for a linked list, but it’s short enough to fit in a book. Please don’t use it in real code.

If you see a function that takes in an empty interface, it’s likely that it is using reflection (which we’ll talk about in Chapter 14) to either populate or read the value. In our preceding example, the second parameter of the json.Unmarshal function is declared to be of type interface{}.

These situations should be relatively rare. Avoid using interface{}. As we’ve seen, Go is designed as a strongly typed language and attempts to work around this are unidiomatic.

If you find yourself in a situation where you had to store a value into an empty interface, you might be wondering how to read the value back again. To do that, we need to look at type assertions and type switches.

Type Assertions and Type Switches
Go provides two ways to see if a variable of an interface type has a specific concrete type or if the concrete type implements another interface. Let’s start by looking at type assertions. A type assertion names the concrete type that implemented the interface, or names another interface that is also implemented by the concrete type underlying the interface. You can try it out on The Go Playground:

type MyInt int

func main() {
    var i interface{}
    var mine MyInt = 20
    i = mine
    i2 := i.(MyInt)
    fmt.Println(i2 + 1)
}
In the preceding code, the variable i2 is of type MyInt.

You might wonder what happens if a type assertion is wrong. In that case, your code panics. You can try it out on The Go Playground:

i2 := i.(string)
fmt.Println(i2)
Running this code produces the following panic:

panic: interface conversion: interface {} is main.MyInt, not string
As we’ve already seen, Go is very careful about concrete types. Even if two types share an underlying type, a type assertion must match the type of the underlying value. The following code panics. You can try it out on The Go Playground:

i2 := i.(int)
fmt.Println(i2 + 1)
Obviously, crashing is not desired behavior. We avoid this by using the comma ok idiom, just as we saw in “The comma ok Idiom” when detecting whether or not a zero value was in a map:

i2, ok := i.(int)
if !ok {
    return fmt.Errorf("unexpected type for %v",i)
}
fmt.Println(i2 + 1)
The boolean ok is set to true if the type conversion was successful. If it was not, ok is set to false and the other variable (in this case i2) is set to its zero value. We then handle the unexpected condition within an if statement, but in idiomatic Go, we indent the error handling code. We’ll talk more about error handling in Chapter 8.

NOTE
A type assertion is very different from a type conversion. Type conversions can be applied to both concrete types and interfaces and are checked at compilation time. Type assertions can only be applied to interface types and are checked at runtime. Because they are checked at runtime, they can fail. Conversions change, assertions reveal.

Even if you are absolutely certain that your type assertion is valid, use the comma ok idiom version. You don’t know how other people (or you in six months) will reuse your code. Sooner or later, your unvalidated type assertions will fail at runtime.

When an interface could be one of multiple possible types, use a type switch instead:

func doThings(i interface{}) {
    switch j := i.(type) {
    case nil:
        // i is nil, type of j is interface{}
    case int:
        // j is of type int
    case MyInt:
        // j is of type MyInt
    case io.Reader:
        // j is of type io.Reader
    case string:
        // j is a string
    case bool, rune:
        // i is either a bool or rune, so j is of type interface{}
    default:
        // no idea what i is, so j is of type interface{}
    }
}
A type switch looks a lot like the switch statement that we saw way back in “switch”. Instead of specifying a boolean operation, you specify a variable of an interface type and follow it with .(type). Usually, you assign the variable being checked to another variable that’s only valid within the switch.

NOTE
Since the purpose of a type switch is to derive a new variable from an existing one, it is idiomatic to assign the variable being switched on to a variable of the same name (i := i.(type)), making this one of the few places where shadowing is a good idea. To make the comments more readable, our example doesn’t use shadowing.

The type of the new variable depends on which case matches. You can use nil for one case to see if the interface has no associated type. If you list more than one type on a case, the new variable is of type interface{}. Just like a switch statement, you can have a default case that matches when no specified type does. Otherwise, the new variable has the type of the case that matches.

TIP
If you don’t know the underlying type, you need to use reflection. We’ll talk more about reflection in Chapter 14.

Use Type Assertions and Type Switches Sparingly
While it might seem handy to be able to extract the concrete implementation from an interface variable, you should use these techniques infrequently. For the most part, treat a parameter or return value as the type that was supplied and not what else it could be. Otherwise, your function’s API isn’t accurately declaring what types it needs to perform its task. If you needed a different type, then it should be specified.

That said, there are use cases where type assertions and type switches are useful. One common use of a type assertion is to see if the concrete type behind the interface also implements another interface. This allows you to specify optional interfaces. For example, the standard library uses this technique to allow more efficient copies when the io.Copy function is called. This function has two parameters of types io.Writer and io.Reader and calls the io.copyBuffer function to do its work. If the io.Writer parameter also implements io.WriterTo, or the io.Reader parameter also implements io.ReaderFrom, most of the work in the function can be skipped:

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
    // If the reader has a WriteTo method, use it to do the copy.
    // Avoids an allocation and a copy.
    if wt, ok := src.(WriterTo); ok {
        return wt.WriteTo(dst)
    }
    // Similarly, if the writer has a ReadFrom method, use it to do the copy.
    if rt, ok := dst.(ReaderFrom); ok {
        return rt.ReadFrom(src)
    }
    // function continues...
}
Another place optional interfaces are used is when evolving an API. In Chapter 12 we’ll discuss the context. Context is a parameter that’s passed to functions that provides, among other things, a standard way to manage cancellation. It was added to Go in version 1.7, which means older code doesn’t support it. This includes older database drivers.

In Go 1.8, new context-aware analogues of existing interfaces were defined in the database/sql/driver package. For example, the StmtExecContext interface defines a method called ExecContext, which is a context-aware replacement for the Exec method in Stmt. When an implementation of Stmt is passed into standard library database code, it checks to see if it also implements StmtExecContext. If it does, ExecContext is invoked. If not, the Go standard library provides a fallback implementation of the cancellation support provided by newer code:

func ctxDriverStmtExec(ctx context.Context, si driver.Stmt,
                       nvdargs []driver.NamedValue) (driver.Result, error) {
    if siCtx, is := si.(driver.StmtExecContext); is {
        return siCtx.ExecContext(ctx, nvdargs)
    }
    // fallback code is here
}
There is one drawback to the optional interface technique. We saw earlier that it is common for implementations of interfaces to use the decorator pattern to wrap other implementations of the same interface to layer behavior. The problem is that if there is an optional interface implemented by one of the wrapped implementations, you cannot detect it with a type assertion or type switch. For example, the standard library includes a bufio package that provides a buffered reader. You can buffer any other io.Reader implementation by passing it to the bufio.NewReader function and using the returned *bufio.Reader. If the passed-in io.Reader also implemented io.ReaderFrom, wrapping it in a buffered reader prevents the optimization.

We also see this when handling errors. As mentioned earlier, they implement the error interface. Errors can include additional information by wrapping other errors. A type switch or type assertion cannot detect or match wrapped errors. If you want different behaviors to handle different concrete implementations of a returned error, use the errors.Is and errors.As functions to test for and access the wrapped error.

Type switch statements provide the ability to differentiate between multiple implementations of an interface that require different processing. They are most useful when there are only certain possible valid types that can be supplied for an interface. Be sure to include a default case in the type switch to handle implementations that aren’t known at development time. This protects you if you forget to update your type switch statements when adding new interface implementations:

func walkTree(t *treeNode) (int, error) {
    switch val := t.val.(type) {
    case nil:
        return 0, errors.New("invalid expression")
    case number:
        // we know that t.val is of type number, so return the
        // int value
        return int(val), nil
    case operator:
        // we know that t.val is of type operator, so
        // find the values of the left and right children, then
        // call the process() method on operator to return the
        // result of processing their values.
        left, err := walkTree(t.lchild)
        if err != nil {
            return 0, err
        }
        right, err := walkTree(t.rchild)
        if err != nil {
            return 0, err
        }
        return val.process(left, right), nil
    default:
        // if a new treeVal type is defined, but walkTree wasn't updated
        // to process it, this detects it
        return 0, errors.New("unknown node type")
    }
}
You can see the complete implementation on The Go Playground.

NOTE
You can further protect yourself from unexpected interface implementations by making the interface unexported and at least one method unexported. If the interface is exported, then it can be embedded in a struct in another package, making the struct implement the interface. We’ll talk more about packages and exporting identifiers in Chapter 9.

Function Types Are a Bridge to Interfaces
There’s one last thing that we haven’t talked about with type declarations. It’s pretty easy to wrap your head around adding a method to an int or a string, but Go allows methods on any user-defined type, including user-defined function types. This sounds like an academic corner case, but they are actually very useful. They allow functions to implement interfaces. The most common usage is for HTTP handlers. An HTTP handler processes an HTTP server request. It’s defined by an interface:

type Handler interface {
    ServeHTTP(http.ResponseWriter, *http.Request)
}
By using a type conversion to http.HandlerFunc, any function that has the signature func(http.ResponseWriter,*http.Request) can be used as an http.Handler:

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    f(w, r)
}
This lets you implement HTTP handlers using functions, methods, or closures using the exact same code path as the one used for other types that meet the http.Handler interface.

Functions in Go are first-class concepts, and as such, they are often passed as parameters into functions. Meanwhile, Go encourages small interfaces, and an interface of only one method could easily replace a parameter of function type. The question becomes: when should your function or method specify an input parameter of a function type and when should you use an interface?

If your single function is likely to depend on many other functions or other state that’s not specified in its input parameters, use an interface parameter and define a function type to bridge a function to the interface. That’s what’s done in the http package; it’s likely that a Handler is just the entry point for a chain of calls that needs to be configured. However, if it’s a simple function (like the one used in sort.Slice), then a parameter of function type is a good choice.

Implicit Interfaces Make Dependency Injection Easier
Anyone who has been programming for any length of time quickly learns that applications need to change over time. One of the techniques that has been developed to ease decoupling is called dependency injection. Dependency injection is the concept that your code should explicitly specify the functionality it needs to perform its task. It’s quite a bit older than you might think; in 1996, Robert Martin wrote an article called “The Dependency Inversion Principle”.

One of the surprising benefits of Go’s implicit interfaces is that they make dependency injection an excellent way to decouple your code. While developers in other languages often use large, complicated frameworks to inject their dependencies, the truth is that it is easy to implement dependency injection in Go without any additional libraries. Let’s work through a simple example to see how we use implicit interfaces to compose applications via dependency injection.

To understand this concept better and see how to implement dependency injection in Go, let’s build a very simple web application. (We’ll talk more about Go’s built-in HTTP server support in “The Server”; consider this a preview.) We’ll start by writing a small utility function, a logger:

func LogOutput(message string) {
    fmt.Println(message)
}
Another thing our app needs is a data store. Let’s create a simple one:

type SimpleDataStore struct {
    userData map[string]string
}

func (sds SimpleDataStore) UserNameForID(userID string) (string, bool) {
    name, ok := sds.userData[userID]
    return name, ok
}
Let’s also define a factory function to create an instance of a SimpleDataStore:

func NewSimpleDataStore() SimpleDataStore {
    return SimpleDataStore{
        userData: map[string]string{
            "1": "Fred",
            "2": "Mary",
            "3": "Pat",
        },
    }
}
Next, we’ll write some business logic that looks up a user and says hello or goodbye. Our business logic needs some data to work with, so it requires a data store. We also want our business logic to log when it is invoked, so it depends on a logger. However, we don’t want to force it to depend on LogOutput or SimpleDataStore, because we might want to use a different logger or data store later. What our business logic needs are interfaces to describe what it depends on:

type DataStore interface {
    UserNameForID(userID string) (string, bool)
}

type Logger interface {
    Log(message string)
}
To make our LogOutput function meet this interface, we define a function type with a method on it:

type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
    lg(message)
}
By a stunning coincidence, our LoggerAdapter and SimpleDataStore happen to meet the interfaces needed by our business logic, but neither type has any idea that it does.

Now that we have the dependencies defined, let’s look at the implementation of our business logic:

type SimpleLogic struct {
    l  Logger
    ds DataStore
}

func (sl SimpleLogic) SayHello(userID string) (string, error) {
    sl.l.Log("in SayHello for " + userID)
    name, ok := sl.ds.UserNameForID(userID)
    if !ok {
        return "", errors.New("unknown user")
    }
    return "Hello, " + name, nil
}

func (sl SimpleLogic) SayGoodbye(userID string) (string, error) {
    sl.l.Log("in SayGoodbye for " + userID)
    name, ok := sl.ds.UserNameForID(userID)
    if !ok {
        return "", errors.New("unknown user")
    }
    return "Goodbye, " + name, nil
}
We have a struct with two fields, one a Logger, the other a DataStore. There’s nothing in SimpleLogic that mentions the concrete types, so there’s no dependency on them. There’s no problem if we later swap in new implementations from an entirely different provider, because the provider has nothing to do with our interface. This is very different from explicit interfaces in languages like Java. Even though Java uses an interface to decouple implementation from interface, the explicit interfaces bind the client and the provider together. This makes replacing a dependency in Java (and other languages with explicit interfaces) far more difficult than it is in Go.

When we want a SimpleLogic instance, we call a factory function, passing in interfaces and returning a struct:

func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
    return SimpleLogic{
        l:    l,
        ds: ds,
    }
}
NOTE
The fields in SimpleLogic are unexported. This means they can only be accessed by code within the same package as SimpleLogic. We can’t enforce immutability in Go, but limiting which code can access these fields makes their accidental modification less likely. We’ll talk more about exported and unexported identifiers in Chapter 9.

Now we get to our API. We’re only going to have a single endpoint, /hello, which says hello to the person whose user ID is supplied. (Please do not use query parameters in your real applications for authentication information; this is just a quick sample.) Our controller needs business logic that says hello, so we define an interface for that:

type Logic interface {
    SayHello(userID string) (string, error)
}
This method is available on our SimpleLogic struct, but once again, the concrete type is not aware of the interface. Furthermore, the other method on SimpleLogic, SayGoodbye, is not in the interface because our controller doesn’t care about it. The interface is owned by the client code, so its method set is customized to the needs of the client code:

type Controller struct {
    l     Logger
    logic Logic
}

func (c Controller) SayHello(w http.ResponseWriter, r *http.Request) {
    c.l.Log("In SayHello")
    userID := r.URL.Query().Get("user_id")
    message, err := c.logic.SayHello(userID)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }
    w.Write([]byte(message))
}
Just as we have factory functions for our other types, let’s write one for the Controller:

func NewController(l Logger, logic Logic) Controller {
    return Controller{
        l:     l,
        logic: logic,
    }
}
Again, we accept interfaces and return structs.

Finally, we wire up all of our components in our main function and start our server:

func main() {
    l := LoggerAdapter(LogOutput)
    ds := NewSimpleDataStore()
    logic := NewSimpleLogic(l, ds)
    c := NewController(l, logic)
    http.HandleFunc("/hello", c.SayHello)
    http.ListenAndServe(":8080", nil)
}
The main function is the only part of the code that knows what all the concrete types actually are. If we want to swap in different implementations, this is the only place that needs to change. Externalizing the dependencies via dependency injection means that we limit the changes that are needed to evolve our code over time.

Dependency injection is also a great pattern for making testing easier. It shouldn’t be surprising, since writing unit tests is effectively reusing your code in a different environment, one where the inputs and outputs are constrained to validate functionality. For example, we can validate the logging output in a test by injecting a type that captures the log output and meets the Logger interface. We’ll talk about this more in Chapter 13.

NOTE
The line http.HandleFunc("/hello", c.SayHello) demonstrates two things we talked about earlier.

First, we are treating the SayHello method as a function.

Second, the http.HandleFunc function takes in a function and converts it to an http.HandlerFunc function type, which declares a method to meet the http.Handler interface, which is the type used to represent a request handler in Go. We took a method from one type and converted it into another type with its own method. That’s pretty neat.