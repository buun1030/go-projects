Servers need a way to handle metadata on individual requests. This metadata falls into two general categories: metadata that is required to correctly process the request, and metadata on when to stop processing the request. For example, an HTTP server might want to use a tracking ID to identify a chain of requests through a set of microservices. It also might want to set a timer that ends requests to other microservices if they take too long. Many languages use threadlocal variables to store this kind of information, associating data to a specific operating system thread of execution. This does’t work in Go because goroutines don’t have unique identities that can be used to look up values. More importantly, threadlocals feel like magic; values go in one place and pop up somewhere else.

Go solves the request metadata problem with a construct called the context. Let’s see how to use it correctly.

What Is the Context?
Rather than add a new feature to the language, a context is simply an instance that meets the Context interface defined in the context package. As you know, idiomatic Go encourages explicit data passing via function parameters. The same is true for the context. It is just another parameter to your function. Just like Go has a convention that the last return value from a function is an error, there is another Go convention that the context is explicitly passed through your program as the first parameter of a function. The usual name for the context parameter is ctx:

func logic(ctx context.Context, info string) (string, error) {
    // do some interesting stuff here
    return "", nil
}
In addition to defining the Context interface, the context package also contains several factory functions for creating and wrapping contexts. When you don’t have an existing context, such as at the entry point to a command-line program, create an empty initial context with the function context.Background. This returns a variable of type context.Context. (Yes, this is an exception to the usual pattern of returning a concrete type from a function call.)

An empty context is a starting point; each time you add metadata to the context, you do so by wrapping the existing context using one of the factory functions in the context package:

ctx := context.Background()
result, err := logic(ctx, "a string")
NOTE
There is another function, context.TODO, that also creates an empty context.Context. It is intended for temporary use during development. If you aren’t sure where the context is going to come from or how it’s going to be used, use context.TODO to put a placeholder in your code. Production code shouldn’t include context.TODO.

When writing an HTTP server, you use a slightly different pattern for acquiring and passing the context through layers of middleware to the top-level http.Handler. Unfortunately, context was added to the Go APIs long after the net/http package was created. Due to the compatibility promise, there was no way to change the http.Handler interface to add a context.Context parameter.

The compatibility promise does allow new methods to be added to existing types, and that’s what the Go team did. There are two context-related methods on http.Request:

Context returns the context.Context associated with the request.

WithContext takes in a context.Context and returns a new http.Request with the old request’s state combined with the supplied context.Context.

Here’s the general pattern:

func Middleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        // wrap the context with stuff -- we'll see how soon!
        req = req.WithContext(ctx)
        handler.ServeHTTP(rw, req)
    })
}
The first thing we do in our middleware is extract the existing context from the request using the Context method. After we put values into the context, we create a new request based on the old request and the now-populated context using the WithContext method. Finally, we call the handler and pass it our new request and the existing http.ResponseWriter.

When you get to the handler, you extract the context from the request using the Context method and call your business logic with the context as the first parameter, just like we saw previously:

func handler(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    err := req.ParseForm()
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    data := req.FormValue("data")
    result, err := logic(ctx, data)
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    rw.Write([]byte(result))
}
There’s one more situation where you use the WithContext method: when making an HTTP call from your application to another HTTP service. Just like we did when passing a context through middleware, you set the context on the outgoing request using WithContext:

type ServiceCaller struct {
    client *http.Client
}

func (sc ServiceCaller) callAnotherService(ctx context.Context, data string)
                                          (string, error) {
    req, err := http.NewRequest(http.MethodGet,
                "http://example.com?data="+data, nil)
    if err != nil {
        return "", err
    }
    req = req.WithContext(ctx)
    resp, err := sc.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("Unexpected status code %d",
                              resp.StatusCode)
    }
    // do the rest of the stuff to process the response
    id, err := processResponse(resp.Body)
    return id, err
}
Now that we know how to acquire and pass a context, let’s start making them useful. We’ll begin with cancellation.

Cancellation
Imagine that you have a request that spawns several goroutines, each one calling a different HTTP service. If one service returns an error that prevents you from returning a valid result, there is no point in continuing to process the other goroutines. In Go, this is called cancellation and the context provides the mechanism for implementation.

To create a cancellable context, use the context.WithCancel function. It takes in a context.Context as a parameter and returns a context.Context and a context.CancelFunc. The returned context.Context is not the same context that was passed into the function. Instead, it is a child context that wraps the passed-in parent context.Context. A context.CancelFunc is a function that cancels the context, telling all of the code that’s listening for potential cancellation that it’s time to stop processing.

NOTE
We’ll see this wrapping pattern several times. A context is treated as an immutable instance. Whenever we add information to a context, we do so by wrapping an existing parent context with a child context. This allows us to use contexts to pass information into deeper layers of the code. The context is never used to pass information out of deeper layers to higher layers.

Let’s take a look at how it works. Because this code sets up a server, you can’t run it on The Go Playground, but you can download it. First we’ll set up two servers in a file called servers.go:

func slowServer() *httptest.Server {
    s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
         r *http.Request) {
        time.Sleep(2 * time.Second)
        w.Write([]byte("Slow response"))
    }))
    return s
}

func fastServer() *httptest.Server {
    s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
         r *http.Request) {
        if r.URL.Query().Get("error") == "true" {
            w.Write([]byte("error"))
            return
        }
        w.Write([]byte("ok"))
    }))
    return s
}
These functions launch servers when they are called. One server sleeps for two seconds and then returns the message Slow response. The other checks to see if there is a query parameter error set to true. If there is, it returns the message error. Otherwise, it returns the message ok.

NOTE
We are using the httptest.Server, which makes it easier to write unit tests for code that talks to remote servers. It’s useful here since both the client and the server are within the same program. We’ll learn more about httptest.Server in “httptest”.

Next, we’re going to write the client portion of the code in a file called client.go:

var client = http.Client{}

func callBoth(ctx context.Context, errVal string, slowURL string,
              fastURL string) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        defer wg.Done()
        err := callServer(ctx, "slow", slowURL)
        if err != nil {
            cancel()
        }
    }()
    go func() {
        defer wg.Done()
        err := callServer(ctx, "fast", fastURL+"?error="+errVal)
        if err != nil {
            cancel()
        }
    }()
    wg.Wait()
    fmt.Println("done with both")
}

func callServer(ctx context.Context, label string, url string) error {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        fmt.Println(label, "request err:", err)
        return err
    }
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(label, "response err:", err)
        return err
    }
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(label, "read err:", err)
        return err
    }
    result := string(data)
    if result != "" {
        fmt.Println(label, "result:", result)
    }
    if result == "error" {
        fmt.Println("cancelling from", label)
        return errors.New("error happened")
    }
    return nil
}
All of the interesting stuff is in this file. First, our callBoth function creates a cancellable context and a cancellation function from the passed-in context. By convention, this function variable is named cancel. It is important to remember that any time you create a cancellable context, you must call the cancel function. It is fine to call it more than once; every invocation after the first is ignored. We use a defer to make sure that it is eventually called. Next, we set up two goroutines and pass the cancellable context, a label, and the URL to callServer, and wait for them both to complete. If either call to callServer returns an error, we call the cancel function.

The callServer function is a simple client. We create our requests with the cancellable context and make a call. If an error happens, or if we get the string error returned, we return the error.

Finally, we have the main function, which kicks off the program, in the file main.go:

func main() {
    ss := slowServer()
    defer ss.Close()
    fs := fastServer()
    defer fs.Close()

    ctx := context.Background()
    callBoth(ctx, os.Args[1], ss.URL, fs.URL)
}
In main, we start the servers, create a context, and then call the clients with the context, the first argument to our program, and the URLs for our servers.

Here’s what happens if you run without an error:

$ make run-ok
go build
./context_cancel false
fast result: ok
slow result: Slow response
done with both
And here’s what happens if an error is triggered:

$ make run-cancel
go build
./context_cancel true
fast result: error
cancelling from fast
slow response err: Get "http://127.0.0.1:38804": context canceled
done with both
NOTE
Any time you create a context that has an associated cancel function, you must call that cancel function when you are done processing, whether or not your processing ends in an error. If you do not, your program will leak resources (memory and goroutines) and eventually slow down or crash. There is no error if you call the cancel function more than once; any invocation after the first does nothing. The easiest way to make sure you call the cancel function is to use defer to invoke it right after the cancel function is returned.

While manual cancellation is useful, it’s not your only option. In the next section, we’ll see how to automate cancellation with timeouts.

HANDLING SERVER SHUTDOWN
You probably noticed that the program didn’t exit immediately when the error is triggered. That’s because it’s waiting for the slowServer to close. If you lengthen the timeout from 2 seconds to 6 seconds or more, you’ll see an error message that starts with httptest.Server blocked in Close after 5 seconds, waiting for connections.

If we rewrite slowServer() to properly handle context cancellation, we can shut it down immediately:

func slowServer() *httptest.Server {
    s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
                                                  r *http.Request) {
        ctx := r.Context()
        select {
        case <-ctx.Done():
            fmt.Println("server shut down")
            return
        case <-time.After(6 * time.Second):
            w.Write([]byte("Slow response"))
        }
    }))
    return s
}
Timers
One of the most important jobs for a server is managing requests. A novice programmer often thinks that a server should take as many requests as it possibly can and work on them for as long as it can until it returns a result for each client.

The problem is that this approach does not scale. A server is a shared resource. Like all shared resources, each user wants to get as much as they can out of it and isn’t terribly concerned with the needs of other users. It’s the responsibility of the shared resource to manage itself so that it provides a fair amount of time to all of its users.

There are generally four things that a server can do to manage its load:

Limit simultaneous requests

Limit how many requests are queued waiting to run

Limit how long a request can run

Limit the resources a request can use (such as memory or disk space)

Go provides tools to handle the first three. We saw how to handle the first two when learning about concurrency in Chapter 10. By limiting the number of goroutines, a server manages simultaneous load. The size of the waiting queue is handled via buffered channels.

The context provides a way to control how long a request runs. When building an application, you should have an idea of your performance envelope: how long you have for your request to complete before the user has an unsatisfactory experience. If you know the maximum amount of time that a request can run, you can enforce it using the context.

NOTE
If you want to limit the memory or disk space that a request uses, you’ll have to write the code to manage that yourself. Discussion of this topic is beyond the scope of this book.

You can use one of two different functions to create a time-limited context. The first is context.WithTimeout. It takes two parameters, an existing context and time.Duration that specifies the duration until the context automatically cancels. It returns a context that automatically triggers a cancellation after the specified duration as well as a cancellation function that is invoked to cancel the context immediately.

The second function is context.WithDeadline. This function takes in an existing context and a time.Time that specifies the time when the context is automatically canceled. Like context.WithTimeout, it returns a context that automatically triggers a cancellation after the specified time has elapsed as well as a cancellation function.

TIP
If you pass a time in the past to context.WithDeadline, the context is created already canceled.

If you want to find out when a context will automatically cancel, use the Deadline method on context.Context. It returns a time.Time that indicates the time and a bool that indicates if there was a timeout set. This mirrors the comma ok idiom we use when reading from maps or channels.

When you set a time limit for the overall duration of the request, you might want to subdivide that time. And if you call another service from your service, you might want to limit how long you allow the network call to run, reserving some time for the rest of your processing or for other network calls. You control how long an individual call takes by creating a child context that wraps a parent context using context.WithTimeout or context.WithDeadline.

Any timeout that you set on the child context is bounded by the timeout set on the parent context; if a parent context times out in two seconds, you can declare that a child context times out in three seconds, but when the parent context times out after two seconds, so will the child.

We can see this with a simple program:

ctx := context.Background()
parent, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()
child, cancel2 := context.WithTimeout(parent, 3*time.Second)
defer cancel2()
start := time.Now()
<-child.Done()
end := time.Now()
fmt.Println(end.Sub(start))
In this sample, we specify a two-second timeout on the parent context and a three-second timeout on the child context. We then wait for the child context to complete by waiting on the channel returned from the Done method on the child context.Context. We’ll talk more about the Done method in the next section.

You can run this code on The Go Playground and you’ll see the following result:

2s
Handling Context Cancellation in Your Own Code
Most of the time, you don’t need to worry about timeouts or cancellation within your own code; it simply doesn’t run for long enough. Whenever you call another HTTP service or the database, you should pass along the context; those libraries properly handle cancellation via the context.

If you do write code that should be interrupted by a context cancellation, you implement the cancellation checks using the concurrency features that we looked at in Chapter 10. The context.Context type has two methods that are used when managing cancellation.

The Done method returns a channel of struct{}. (The reason this is the chosen return type is that an empty struct uses no memory.) The channel is closed when the context is canceled due to a timer or the cancel function being invoked. Remember, a closed channel always immediately returns its zero value when you attempt to read it.

WARNING
If you call Done on a context that isn’t cancellable, it returns nil. As we covered in Chapter 10, a read from a nil channel never returns. If this is not done inside a case in a select statement, your program will hang.

The Err method returns nil if the context is still active, or it returns one of two sentinel errors if the context has been canceled: context.Canceled and context.DeadlineExceeded. The first is returned after explicit cancellation, and the second is returned when a timeout triggered cancellation.

Here’s the pattern for supporting context cancellation in your code:

func longRunningThingManager(ctx context.Context, data string) (string, error) {
    type wrapper struct {
        result string
        err    error
    }
    ch := make(chan wrapper, 1)
    go func() {
        // do the long running thing
        result, err := longRunningThing(ctx, data)
        ch <- wrapper{result, err}
    }()
    select {
    case data := <-ch:
        return data.result, data.err
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
In our code, we need to put the data returned from our long-running function into a struct, so we can pass it on a channel. We then create a channel of type wrapper with buffer size 1. By buffering the channel, we allow the goroutine to exit, even if the buffered value is never read due to cancellation.

In the goroutine, we take the output from the long-running function and put it in the buffered channel. We then have a select with two cases. In our first select case, we read the data from the long-running function and return it. This is the case that’s triggered if the context isn’t canceled due to timeout or invocation of the cancel function. The second select case is triggered if the context is canceled. We return the zero value for the data and the error from the context to tell us why it was canceled.

This looks a lot like the pattern we saw in Chapter 11, when we learned how to use time.After to set a time limit on the execution of code. In this case, the time limit (or the cancellation condition) is specified via context factory methods, but the general implementation is the same.

Values
There is one more use for the context. It also provides a way to pass per-request metadata through your program.

By default, you should prefer to pass data through explicit parameters. As has been mentioned before, idiomatic Go favors the explicit over the implicit, and this includes explicit data passing. If a function depends on some data, it should be clear where it came from.

However, there are some cases where you cannot pass data explicitly. The most common situation is an HTTP request handler and its associated middleware. As we have seen, all HTTP request handlers have two parameters, one for the request and one for the response. If you want to make a value available to your handler in middleware, you need to store it in the context. Some possible situations include extracting a user from a JWT (JSON Web Token) or creating a per-request GUID that is passed through multiple layers of middleware and into your handler and business logic.

Just like there are factory methods in the context package to create timed and cancellable contexts, there is a factory method for putting values into the context, context.WithValue. It takes in three values: a context to wrap, a key to look up the value, and the value itself. It returns a child context that contains the key-value pair. The type of the key and the value parameters are declared to be empty interfaces (interface{}).

To check if a value is in a context or any of its parents, use the Value method on context.Context. This method takes in a key and returns the value associated with the key. Again, both the key parameter and the value result are declared to be of type interface{}. If no value is found for the supplied key, nil is returned. Use the comma ok idiom to type assert the returned value to the correct type.

NOTE
If you are familiar with data structures, you might recognize that searching for values stored in the context chain is a linear search. This has no serious performance implications when there are only a few values, but it would perform poorly if you stored dozens of values in the context during a request. That said, if your program is creating a context chain with dozens of values, your program probably needs some refactoring.

While the value stored in the context can be of any type, there is an idiomatic pattern that’s used to guarantee the key’s uniqueness. Like the key for a map, the key for context value must be comparable. Create a new, unexported type for the key, based on an int:

type userKey int
If you use a string or another public type for the type of the key, different packages could create identical keys, resulting in collisions. This causes problems that are hard to debug, such as one package writing data to the context that masks the data written by another package or reading data from the context that was written by another package.

After declaring your unexported key type, you then declare an unexported constant of that type:

const key userKey = 1
With both the type and the constant of the key being unexported, no code from outside of your package can put data into the context that would cause a collision. If your package needs to put multiple values into the context, define a different key of the same type for each value, using the iota pattern we looked at in “iota Is for Enumerations—Sometimes”. Since we only care about the constant’s value as a way to differentiate between multiple keys, this is a perfect use for iota.

Next, build an API to place a value into the context and to read the value from the context. Make these functions public only if code outside your package should be able to read and write your context values. The name of the function that creates a context with the value should start with ContextWith. The function that returns the value from the context should have a name that ends with FromContext. Here are the implementations of our functions to get and read the user from the context:

func ContextWithUser(ctx context.Context, user string) context.Context {
    return context.WithValue(ctx, key, user)
}

func UserFromContext(ctx context.Context) (string, bool) {
    user, ok := ctx.Value(key).(string)
    return user, ok
}
Now that we’ve written our user management code, let’s see how to use it. We’re going to write middleware that extracts a user ID from a cookie:

// a real implementation would be signed to make sure
// the user didn't spoof their identity
func extractUser(req *http.Request) (string, error) {
    userCookie, err := req.Cookie("user")
    if err != nil {
        return "", err
    }
    return userCookie.Value, nil
}

func Middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        user, err := extractUser(req)
        if err != nil {
            rw.WriteHeader(http.StatusUnauthorized)
            return
        }
        ctx := req.Context()
        ctx = ContextWithUser(ctx, user)
        req = req.WithContext(ctx)
        h.ServeHTTP(rw, req)
    })
}
In the middleware, we first get our user value. Next, we extract the context from the request with the Context method and create a new context that contains the user with our ContextWithUser function. Then we make a new request from the old request and the new context using the WithContext method. Finally, we call the next function in our handler chain with our new request and the supplied http.ResponseWriter.

In most cases, you want to extract the value from the context in your request handler and pass it in to your business logic explicitly. Go functions have explicit parameters and you shouldn’t use the context as a way to sneak values past the API:

func (c Controller) handleRequest(rw http.ResponseWriter, req *http.Request) {
    ctx := req.Context()
    user, ok := identity.UserFromContext(ctx)
    if !ok {
        rw.WriteHeader(http.StatusInternalServerError)
        return
    }
    data := req.URL.Query().Get("data")
    result, err := c.Logic.businessLogic(ctx, user, data)
    if err != nil {
        rw.WriteHeader(http.StatusInternalServerError)
        rw.Write([]byte(err.Error()))
        return
    }
    rw.Write([]byte(result))
}
Our handler gets the context using the Context method on the request, extracts the user from the context using our UserFromContext function, and then calls the business logic.

There are some situations where it’s better to keep a value in the context. The tracking GUID that was mentioned earlier is one. This information is meant for management of your application; it is not part of your business state. Passing it explicitly through your code adds additional parameters and prevents integration with third-party libraries that do not know about your metainformation. By leaving a tracking GUID in the context, it passes invisibly through business logic that doesn’t need to know about tracking and is available when your program writes a log message or connects to another server.

Here is a simple context-aware GUID implementation that tracks from service to service and creates logs with the GUID included:

package tracker

import (
    "context"
    "fmt"
    "net/http"
    "github.com/google/uuid"
)

type guidKey int

const key guidKey = 1

func contextWithGUID(ctx context.Context, guid string) context.Context {
    return context.WithValue(ctx, key, guid)
}

func guidFromContext(ctx context.Context) (string, bool) {
    g, ok := ctx.Value(key).(string)
    return g, ok
}

func Middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx := req.Context()
        if guid := req.Header.Get("X-GUID"); guid != "" {
            ctx = contextWithGUID(ctx, guid)
        } else {
            ctx = contextWithGUID(ctx, uuid.New().String())
        }
        req = req.WithContext(ctx)
        h.ServeHTTP(rw, req)
    })
}

type Logger struct{}

func (Logger) Log(ctx context.Context, message string) {
    if guid, ok := guidFromContext(ctx); ok {
        message = fmt.Sprintf("GUID: %s - %s", guid, message)
    }
    // do logging
    fmt.Println(message)
}

func Request(req *http.Request) *http.Request {
    ctx := req.Context()
    if guid, ok := guidFromContext(ctx); ok {
        req.Header.Add("X-GUID", guid)
    }
    return req
}
The Middleware function either extracts the GUID from the incoming request or generates a new GUID. In both cases, it places the GUID into the context, creates a new request with the updated context, and continues the call chain.

Next we see how this GUID is used. The Logger struct provides a generic logging method that takes in a context and a string. If there’s a GUID in the context, it appends it to the beginning of the log message and outputs it. The Request function is used when this service makes a call to another service. It takes in an *http.Request, adds a header with the GUID if it exists in the context, and returns the *http.Request.

Once we have this package, we can use the dependency injection techniques that we discussed in “Implicit Interfaces Make Dependency Injection Easier” to create business logic that is completely unaware of any tracking information. First, we declare an interface to represent our logger, a function type to represent a request decorator, and a business logic struct that depends on them:

type Logger interface {
    Log(context.Context, string)
}

type RequestDecorator func(*http.Request) *http.Request

type BusinessLogic struct {
    RequestDecorator RequestDecorator
    Logger                     Logger
    Remote                     string
}
Next, we implement our business logic:

func (bl BusinessLogic) businessLogic(
    ctx context.Context, user string, data string) (string, error) {
    bl.Logger.Log(ctx, "starting businessLogic for " + user + " with "+ data)
    req, err := http.NewRequestWithContext(ctx,
        http.MethodGet, bl.Remote+"?query="+data, nil)
    if err != nil {
        bl.Logger.Log(ctx, "error building remote request:" + err)
        return "", err
    }
    req = bl.RequestDecorator(req)
    resp, err := http.DefaultClient.Do(req)
    // processing continues
}
The GUID is passed through to the logger and the request decorator without the business logic being aware of it, separating the data needed for program logic from the data needed for program management. The only place that’s aware of the association is the code in main that wires up our dependencies:

bl := BusinessLogic{
    RequestDecorator: tracker.Request,
    Logger:           tracker.Logger{},
    Remote:           "http://www.example.com/query",
}
You can find the complete code for the user middleware and the GUID tracker on GitHub.

TIP
Use the context to pass values through standard APIs. Copy values from the context into explicit parameters when they are needed for processing business logic. System maintenance information can be accessed directly from the context.