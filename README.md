# go-mistakes

## Unintended variable shadowing

Declares outer block variable and then redeclare the same name variable in an inner block,
after prcess the inner block, the outer block variable will stay the same.\
\
**Solution:** To assign value to the outer block variable in the inner block, just use `=` not `:=`

## Unnecessary nested code

Because it was difficult to distinguish the expected execution flow because of the nested if/else statements. Conversely, if only one nested if/else statment, it requires scanning down one column to see the expected execution flow and down the second column to see how the edge cases are handled.\
\
![This is an alt text.](![Alt text](image-1.png) "To understand the expected execution flow, we just have to scan the happy path column.")\
\
**Solution:** Striving to reduce the number of nested blocks, aligning the happy path on the left, and returning as early as possible are concrete means to improve our code’s readability.

## Interface on the producer side

*abstractions should be discovered, not created.* This means that it’s not up to the producer to force a given abstraction for all the clients. Instead, it’s up to the client to decide whether it needs some form of abstraction and then determine the best abstraction level for its needs.\
\
`time.Time` field contain monotonic time An interface should live on the consumer side in most cases. However, in particular contexts (for example, when we know—not foresee—that an abstraction will be helpful for consumers), we may want to have it on the producer side. If we do, we should strive to keep it as minimal as possible, increasing its reusability potential and making it more easily composable.\
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

## Common JSON-handling mistakes

### Unexpected behavior due to type embedding

 If an embedded field type implements an interface, the struct containing the embedded field will also implement this interface. We should be careful with embedded fields. While promoting the fields and methods of an embedded field type can sometimes be convenient, it can also lead to subtle bugs because it can make the parent struct implement interfaces without a clear signal. Again, when using embedded fields

 ### JSON and the monotonic clock

```
2021-01-10 17:13:08.852061 +0100 CET m=+0.000338660
------------------------------------ --------------
             Wall time               Monotonic time
```

When marshal `time.Time` field contain monotonic time while unmarshal `time.Time` field doesn't contain monotonic time. The marshaling/unmarshaling process isn’t always symmetric, and we faced this case with a struct containing a time.Time.\
\
**Solution:** Use the `Equal` method instead of `==` or stay `==` but use `Truncate` method to time before marshal.

### Map of any

```
b := getMessage()
var m map[string]any
err := json.Unmarshal(b, &m)
if err != nil {
    return err
}
```

If we use a map of any: any numeric value, regardless of whether it contains a decimal, is converted into a float64 type.

## Common SQL mistakes

The `database/sql` package provides a generic interface around SQL (or SQL-like) databases. It’s also fairly common to see some patterns or mistakes while using this package. Let’s delve into five common mistakes.

###  Forgetting that sql.Open doesn’t necessarily establish connections to a database

When using sql.Open, one common misconception is expecting this function to establish connections to a database:

```
db, err := sql.Open("mysql", dsn)
if err != nil {
    return err
}
```

But this isn’t necessarily the case. According to the documentation (https://pkg.go.dev/database/sql),
\
*Open may just validate its arguments without creating a connection to the database.*\
\
Actually, the behavior depends on the SQL driver used. For some drivers, sql.Open doesn’t establish a connection: it’s only a preparation for later use (for example, with db.Query). Therefore, the first connection to the database may be established lazily.\
\
Why do we need to know about this behavior? For example, in some cases, we want to make a service ready only after we know that all the dependencies are correctly set up and reachable. If we don’t know this, the service may accept traffic despite an erroneous configuration.\
\
If we want to ensure that the function that uses sql.Open also guarantees that the underlying database is reachable, we should use the Ping method:

```
db, err := sql.Open("mysql", dsn)
if err != nil {
    return err
}
if err := db.Ping(); err != nil {     ❶
    return err
}
```

❶ Calls the Ping method following sql.Open\
\
Ping forces the code to establish a connection that ensures that the data source name is valid and the database is reachable. Note that an alternative to Ping is PingContext, which asks for an additional context conveying when the ping should be canceled or time out.\
\
Despite being perhaps counterintuitive, let’s remember that sql.Open doesn’t necessarily establish a connection, and the first connection can be opened lazily. If we want to test our configuration and be sure a database is reachable, we should follow sql.Open with a call to the Ping or PingContext method.

### Forgetting about connections pooling

Just as the default HTTP client and server provide default behaviors that may not be effective in production (see mistake #81, “Using the default HTTP client and server”), it’s essential to understand how database connections are handled in Go. sql.Open returns an *sql.DB struct. This struct doesn’t represent a single database connection; instead, it represents a pool of connections. This is worth noting so we’re not tempted to implement it manually. A connection in the pool can have two states:

* Already used (for example, by another goroutine that triggers a query)
* Idle (already created but not in use for the time being)

It’s also important to remember that creating a pool leads to four available config parameters that we may want to override. Each of these parameters is an exported method of `*sql.DB`:

* SetMaxOpenConns—Maximum number of open connections to the database (default value: unlimited)
* SetMaxIdleConns—Maximum number of idle connections (default value: 2)
* SetConnMaxIdleTime—Maximum amount of time a connection can be idle before it’s closed (default value: unlimited)
* SetConnMaxLifetime—Maximum amount of time a connection can be held open before it’s closed (default value: unlimited)

Figure shows an example with a maximum of five connections. It has four ongoing connections: three idle and one in use. Therefore, one slot remains available for an extra connection. If a new query comes in, it will pick one of the idle connections (if still available). If there are no more idle connections, the pool will create a new connection if an extra slot is available; otherwise, it will wait until a connection is available.\
\
![Alt text](image.png)\
\
Figure: A connection pool with five connections\
\
So, why should we tweak these config parameters?

* Setting SetMaxOpenConns is important for production-grade applications. Because the default value is unlimited, we should set it to make sure it fits what the underlying database can handle.
* The value of SetMaxIdleConns (default: 2) should be increased if our application generates a significant number of concurrent requests. Otherwise, the application may experience frequent reconnects.
* Setting SetConnMaxIdleTime is important if our application may face a burst of requests. When the application returns to a more peaceful state, we want to make sure the connections created are eventually released.
* Setting SetConnMaxLifetime can be helpful if, for example, we connect to a load-balanced database server. In that case, we want to ensure that our application never uses a connection for too long.

For production-grade applications, we must consider these four parameters. We can also use multiple connection pools if an application faces different use cases.

### Not using prepared statements

A prepared statement is a feature implemented by many SQL databases to execute a repeated SQL statement. Internally, the SQL statement is precompiled and separated from the data provided. There are two main benefits:

* Efficiency—The statement doesn’t have to be recompiled (compilation means parsing + optimization + translation).
* Security—This approach reduces the risks of SQL injection attacks.

Therefore, if a statement is repeated, we should use prepared statements. We should also use prepared statements in untrusted contexts (such as exposing an endpoint on the internet, where the request is mapped to an SQL statement).\
\
To use prepared statements, instead of calling the Query method of *sql.DB, we call Prepare:

```
stmt, err := db.Prepare("SELECT * FROM ORDER WHERE ID = ?")   ❶
if err != nil {
    return err
}
rows, err := stmt.Query(id)                                   ❷
// ...
```

❶ Prepares the statement\
\
❷ Executes the prepared query\
\
We prepare the statement and then execute it while providing the arguments. The first output of the Prepare method is an *sql.Stmt, which can be reused and run concurrently. When the statement is no longer needed, it must be closed using the Close() method.\
\
**NOTE** The Prepare and Query methods have alternatives to provide an additional context: PrepareContext and QueryContext.\
\
For efficiency and security, we need to remember to use prepared statements when it makes sense.
## Not closing transient resources

## Forgetting the return statement after replying to an HTTP request

## Using the default HTTP client and server

## Not using testing utility packages

Unaware of these packages and trying to reinvent the wheel or rely on other solutions that aren’t as handy.

### The httptest package

The [httptest](https://pkg.go.dev/net/http/httptest) package provides utilities for HTTP testing for both clients and servers. Let’s look at these two use cases.\
\
First, let’s see how `httptest` can help us while writing an HTTP server. We will implement a handler that performs some basic actions: writing a header and body, and returning a specific status code. For the sake of clarity, we will omit error handling:

```
func Handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("X-API-VERSION", "1.0")
    b, _ := io.ReadAll(r.Body)
    _, _ = w.Write(append([]byte("hello "), b...))     ❶
    w.WriteHeader(http.StatusCreated)
}
```

❶ Concatenates hello with the request body\
\
An HTTP handler accepts two arguments: the request and a way to write the response. The `httptest` package provides utilities for both. For the request, we can use `httptest.NewRequest` to build an `*http.Request` using an HTTP method, a URL, and a body. For the response, we can use `httptest.NewRecorder` to record the mutations made within the handler. Let’s write a unit test of this handler:

```
func TestHandler(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "http://localhost",     ❶
        strings.NewReader("foo"))
    w := httptest.NewRecorder()                                        ❷
    Handler(w, req)                                                    ❸
 
    if got := w.Result().Header.Get("X-API-VERSION"); got != "1.0" {   ❹
        t.Errorf("api version: expected 1.0, got %s", got)
    }
 
    body, _ := ioutil.ReadAll(wordy)                                   ❺
    if got := string(body); got != "hello foo" {
        t.Errorf("body: expected hello foo, got %s", got)
    }
 
    if http.StatusOK != w.Result().StatusCode {                        ❻
        t.FailNow()
    }
}
```

❶ Builds the request\
\
❷ Creates the response recorder\
\
❸ Calls the handler\
\
❹ Verifies the HTTP header\
\
❺ Verifies the HTTP body\
\
❻ Verifies the HTTP status code\
\
Testing a handler using `httptest` doesn’t test the transport (the HTTP part). The focus of the test is calling the handler directly with a request and a way to record the response. Then, using the response recorder, we write the assertions to verify the HTTP header, body, and status code.\
\
Let’s look at the other side of the coin: testing an HTTP client. We will write a client in charge to query an HTTP endpoint that calculates how long it takes to drive from one coordinate to another. The client looks like this:

```
func (c DurationClient) GetDuration(url string,
    lat1, lng1, lat2, lng2 float64) (
    time.Duration, error) {
    resp, err := c.client.Post(
        url, "application/json",
        buildRequestBody(lat1, lng1, lat2, lng2),
    )
    if err != nil {
        return 0, err
    }
 
    return parseResponseBody(resp.Body)
}
```

This code performs an HTTP POST request to the provided URL and returns the parsed response (let’s say, some JSON).\
\
What if we want to test this client? One option is to use Docker and spin up a mock server to return some preregistered responses. However, this approach makes the test slow to execute. The other option is to use httptest.NewServer to create a local HTTP server based on a handler that we will provide. Once the server is up and running, we can pass its URL to `GetDuration`:

```
func TestDurationClientGet(t *testing.T) {
    srv := httptest.NewServer(                                             ❶
        http.HandlerFunc(
            func(w http.ResponseWriter, r *http.Request) {
                _, _ = w.Write([]byte(`{"duration": 314}`))                ❷
            },
        ),
    )
    defer srv.Close()                                                      ❸
 
    client := NewDurationClient()
    duration, err :=
        client.GetDuration(srv.URL, 51.551261, -0.1221146, 51.57, -0.13)   ❹
    if err != nil {
        t.Fatal(err)
    }
 
    if duration != 314*time.Second {                                       ❺
        t.Errorf("expected 314 seconds, got %v", duration)
    }
}
```

❶ Starts the HTTP server\
\
❷ Registers the handler to serve the response\
\
❸ Shuts down the server\
\
❹ Provides the server URL\
\
❺ Verifies the response\
\
In this test, we create a server with a static handler returning `314` seconds. We could also make assertions based on the request sent. Furthermore, when we call `GetDuration`, we provide the URL of the server that’s started. Compared to testing a handler, this test performs an actual HTTP call, but it executes in only a few milliseconds.\
\
We can also start a new server using TLS with `httptest.NewTLSServer` and create an unstarted server with `httptest.NewUnstartedServer` so that we can start it lazily.\
\
Let’s remember how helpful `httptest` is when working in the context of HTTP applications. Whether we’re writing a server or a client, `httptest` can help us create efficient tests.

### The iotest package

The [iotest](https://pkg.go.dev/testing/iotest) package implements utilities for testing readers and writers. It’s a convenient package that Go developers too often forget.\
\
When implementing a custom `io.Reader`, we should remember to test it using `iotest.TestReader`. This utility function tests that a reader behaves correctly: it accurately returns the number of bytes read, fills the provided slice, and so on. It also tests different behaviors if the provided reader implements interfaces such as `io.ReaderAt`.\
\
Let’s assume we have a custom `LowerCaseReader` that streams lowercase letters from a given input `io.Reader`. Here’s how to test that this reader doesn’t misbehave:

```
func TestLowerCaseReader(t *testing.T) {
    err := iotest.TestReader(
        &LowerCaseReader{reader: strings.NewReader("aBcDeFgHiJ")},   ❶
        []byte("acegi"),                                             ❷
    )
    if err != nil {
        t.Fatal(err)
    }
}
```

❶ Provides an io.Reader\
\
❷ Expectation\
\
We call `iotest.TestReader` by providing the custom `LowerCaseReader` and an expectation: the lowercase letters `acegi`.\
\
Another use case for the `iotest` package is to make sure an application using readers and writers is tolerant to errors:

* `iotest.ErrReader` creates an `io.Reader` that returns a provided error.
* `iotest.HalfReader` creates an `io.Reader` that reads only half as many bytes as requested from an `io.Reader`.
* `iotest.OneByteReader` creates an `io.Reader` that reads a single byte for each non-empty read from an `io.Reader`.
* `iotest.TimeoutReader` creates an `io.Reader` that returns an error on the second read with no data. Subsequent calls will succeed.
* `iotest.TruncateWriter` creates an `io.Writer` that writes to an `io.Writer` but stops silently after n bytes.
\
For example, let’s assume we implement the following function that starts by reading all the bytes from a reader:

```
func foo(r io.Reader) error {
    b, err := io.ReadAll(r)
    if err != nil {
        return err
    }
 
    // ...
}
```

We want to make sure our function is resilient if, for example, the provided reader fails during a read (such as to simulate a network error):

```
func TestFoo(t *testing.T) {
    err := foo(iotest.TimeoutReader(            ❶
        strings.NewReader(randomString(1024)),
    ))
    if err != nil {
        t.Fatal(err)
    }
}
```

❶ Wraps the provided io.Reader using io.TimeoutReader\
\
We wrap an `io.Reader` using `io.TimeoutReader`. As we mentioned, the second read will fail. If we run this test to make sure our function is tolerant to error, we get a test failure. Indeed, `io.ReadAll` returns any errors it finds.\
\
Knowing this, we can implement our custom `readAll` function that tolerates up to n errors:

```
func readAll(r io.Reader, retries int) ([]byte, error) {
    b := make([]byte, 0, 512)
    for {
        if len(b) == cap(b) {
            b = append(b, 0)[:len(b)]
        }
        n, err := r.Read(b[len(b):cap(b)])
        b = b[:len(b)+n]
        if err != nil {
            if err == io.EOF {
                return b, nil
            }
            retries--
            if retries < 0 {     ❶
                return b, err
            }
        }
    }
}
```

❶ Tolerates retries\
\
This implementation is similar to `io.ReadAll`, but it also handles configurable retries. If we change the implementation of our initial function to use our custom `readAll` instead of `io.ReadAll`, the test will no longer fail:

```
func foo(r io.Reader) error {
    b, err := readAll(r, 3)       ❶
    if err != nil {
        return err
    }
 
    // ...
}
```

❶ Indicates up to three retries\
\
We have seen an example of how to check that a function is tolerant to errors while reading from an `io.Reader`. We performed the test by relying on the `iotest` package.\
\
When doing I/O and working with `io.Reader` and `io.Writer`, let’s remember how handy the `iotest` package is. As we have seen, it provides utilities to test the behavior of a custom `io.Reader` and test our application against errors that occur while reading or writing data.

## Not exploring all the go testing features

When it comes to writing tests, developers should know about Go’s specific testing features and options. Otherwise, the testing process can be less accurate and even less efficient. This section discusses topics that can make us more comfortable while writing Go tests.

### Code coverage

During the development process, it can be handy to see visually which parts of our code are covered by tests. We can access this information using the `-coverprofile` flag:\
\
`$ go test -coverprofile=coverage.out ./...`\
\
This command creates a coverage.out file that we can then open using `go tool cover`:\
\
`$ go tool cover -html=coverage.out`\
\
This command opens the web browser and shows the coverage for each line of code.\
\
By default, the code coverage is analyzed only for the current package being tested. For example, suppose we have the following structure:

```
/myapp
  |_ foo
    |_ foo.go
    |_ foo_test.go
  |_ bar
    |_ bar.go
    |_ bar_test.go
```

If some portion of foo.go is only tested in bar_test.go, by default, it won’t be shown in the coverage report. To include it, we have to be in the `myapp` folder and use the `-coverpkg` flag:\
\
`go test -coverpkg=./... -coverprofile=coverage.out ./...`\
\
We need to remember this feature to see the current code coverage and decide which parts deserve more tests.\
\
**NOTE** Remain cautious when it comes to chasing code coverage. Having 100% test coverage doesn’t imply a bug-free application. Properly reasoning about what our tests cover is more important than any static threshold.

### Testing from a different package

When writing unit tests, one approach is to focus on behaviors instead of internals. Suppose we expose an API to clients. We may want our tests to focus on what’s visible from the outside, not the implementation details. This way, if the implementation changes (for example, if we refactor one function into two), the tests will remain the same. They can also be easier to understand because they show how our API is used. If we want to enforce this practice, we can do so using a different package.\
\
In Go, all the files in a folder should belong to the same package, with only one exception: a test file can belong to a `_test` package. For example, suppose the following counter.go source file belongs to the `counter` package:

```
package counter
 
import "sync/atomic"
 
var count uint64
 
func Inc() uint64 {
    atomic.AddUint64(&count, 1)
    return count
}
```

The test file can live in the same package and access internals such as the `count` variable. Or it can live in a `counter_test` package, like this counter_test.go file:

```
package counter_test
 
import (
    "testing"
 
    "myapp/counter"
)
 
func TestCount(t *testing.T) {
    if counter.Inc() != 1 {
        t.Errorf("expected 1")
    }
}
```

In this case, the test is implemented in an external package and cannot access internals such as the `count` variable. Using this practice, we can guarantee that a test won’t use any unexported elements; hence, it will focus on testing the exposed behavior.

### Utility functions

When writing tests, we can handle errors differently than we do in our production code. For example, let’s say we want to test a function that takes as an argument a `Customer` struct. Because the creation of a `Customer` will be reused, we decide to create a specific `createCustomer` function for the sake of the tests. This function will return a possible error alongside a `Customer`:

```
func TestCustomer(t *testing.T) {
    customer, err := createCustomer("foo")     ❶
    if err != nil {
        t.Fatal(err)
    }
    // ...
}
 
func createCustomer(someArg string) (Customer, error) {
    // Create customer
    if err != nil {
        return Customer{}, err
    }
    return customer, nil
}
```

❶ Creates a customer and checks for errors\
\
We create a customer using the `createCustomer` utility function, and then we perform the rest of the test. However, in the context of testing functions, we can simplify error management by passing the `*testing.T` variable to the utility function:

```
func TestCustomer(t *testing.T) {
    customer := createCustomer(t, "foo")     ❶
    // ...
}
 
func createCustomer(t *testing.T, someArg string) Customer {
    // Create customer
    if err != nil {
        t.Fatal(err)                         ❷
    }
    return customer
}
```

❶ Calls the utility function and provides t\
\
❷ Fails the test directly if we can’t create a customer\
\
Instead of returning an error, `createCustomer` fails the test directly if it can’t create a `Customer`. This makes `TestCustomer` smaller to write and easier to read.\
\
Let’s remember this practice regarding error management and testing to improve our tests.

### Setup and teardown

In some cases, we may have to prepare a testing environment. For example, in integration tests, we spin up a specific Docker container and then stop it. We can call setup and teardown functions per test or per package. Fortunately, in Go, both are possible.\
\
To do so per test, we can call the setup function as a preaction and the teardown function using `defer`:

```
func TestMySQLIntegration(t *testing.T) {
    setupMySQL()
    defer teardownMySQL()
    // ...
}
```

It’s also possible to register a function to be executed at the end of a test. For example, let’s assume `TestMySQLIntegration` needs to call `createConnection` to create the database connection. If we want this function to also include the teardown part, we can use `t.Cleanup` to register a cleanup function:

```
func TestMySQLIntegration(t *testing.T) {
    // ...
    db := createConnection(t, "tcp(localhost:3306)/db")
    // ...
}
 
func createConnection(t *testing.T, dsn string) *sql.DB {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        t.FailNow()
    }
    t.Cleanup(          ❶
        func() {
            _ = db.Close()
        })
    return db
}
```

❶ Registers a function to be executed at the end of the test\
\
At the end of the test, the closure provided to `t.Cleanup` is executed. This makes future unit tests easier to write because they won’t be responsible for closing the `db` variable.\
\
Note that we can register multiple cleanup functions. In that case, they will be executed just as if we were using `defer`: last in, first out.\
\
To handle setup and teardown per package, we have to use the `TestMain` function. A simple implementation of `TestMain` is the following:

```
func TestMain(m *testing.M) {
    os.Exit(m.Run())
}
```

This particular function accepts a `*testing.M` argument that exposes a single Run method to run all the tests. Therefore, we can surround this call with setup and teardown functions:

```
func TestMain(m *testing.M) {
    setupMySQL()                 ❶
    code := m.Run()              ❷
    teardownMySQL()              ❸
    os.Exit(code)
}
```

❶ Sets up MySQL\
\
❷ Runs the tests\
\
❸ Tears down MySQL\
\
This code spins up MySQL once before all the tests and then tears it down.\
\
Using these practices to add setup and teardown functions, we can configure a complex environment for our tests.
