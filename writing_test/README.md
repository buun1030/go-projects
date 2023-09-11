Over the past two decades, the widespread adoption of automated testing has probably done more to improve code quality than any other software engineering technique. As a language and ecosystem focused on improving software quality, it’s not surprising that Go includes testing support as part of its standard library. Go makes it so easy to test your code, there’s no excuse to not do it. In this chapter, we’ll see how to test our Go code, group our tests into unit and integration tests, examine code coverage, write benchmarks, and learn how to check our code for concurrency issues using the Go race checker. Along the way, we’ll discuss how to write code that is testable and why this improves our code quality.

The Basics of Testing
Go’s testing support has two parts: libraries and tooling. The testing package in the standard library provides the types and functions to write tests, while the go test tool that’s bundled with Go runs your tests and generates reports. Unlike many other languages, Go tests are placed in the same directory and the same package as the production code. Since tests are located in the same package, they are able to access and test unexported functions and variables. We’ll see in a bit how to write tests that ensure that we are only testing a public API.

NOTE
Complete code samples for this chapter are found on GitHub.

Let’s write a simple function and then a test to make sure the function works. In the file adder/adder.go, we have:

func addNumbers(x, y int) int {
    return x + x
}
The corresponding test is in adder/adder_test.go:

func Test_addNumbers(t *testing.T) {
    result := addNumbers(2,3)
    if result != 5 {
        t.Error("incorrect result: expected 5, got", result)
    }
}
Every test is written in a file whose name ends with _test.go. If you are writing tests against foo.go, place your tests in a file named foo_test.go.

Test functions start with the word Test and take in a single parameter of type *testing.T. By convention, this parameter is named t. Test functions do not return any values. The name of the test (apart from starting with the word “Test”) is meant to document what you are testing, so pick something that explains what you are testing. When writing unit tests for individual functions, the convention is to name the unit test Test followed by the name of the function. When testing unexported functions, some people use an underscore between the word Test and the name of the function.

Also note that we use standard Go code to call the code being tested and to validate if the responses are as expected. When there’s an incorrect result, we report the error with the t.Error method, which works like the fmt.Print function. We’ll see other error-reporting methods in a bit.

We’ve just seen the library portion of Go’s test support. Now let’s take a look at the tooling. Just as go build builds a binary and go run runs a file, the command go test runs the tests in the current directory:

$ go test
--- FAIL: Test_addNumbers (0.00s)
    adder_test.go:8: incorrect result: expected 5, got 4
FAIL
exit status 1
FAIL    test_examples/adder     0.006s
It looks like we found a bug in our code. Taking a second look at addNumbers, we see that we are adding x to x, not x to y. Let’s change the code and rerun our test to verify that the bug is fixed:

$ go test
PASS
ok      test_examples/adder     0.006s
The go test command allows you to specify which packages to test. Using ./… for the package name specifies that you want to run tests in the current directory and all of the subdirectories of the current directory. Include a -v flag to get verbose testing output.

Reporting Test Failures
There are several methods on *testing.T for reporting test failures. We’ve already seen Error, which builds a failure description string out of a comma-separated list of values.

If you’d rather use a Printf-style formatting string to generate your message, use the Errorf method instead:

t.Errorf("incorrect result: expected %d, got %d", 5, result)
While Error and Errorf mark a test as failed, the test function continues running. If you think a test function should stop processing as soon as a failure is found, use the Fatal and Fatalf methods. The Fatal method works like Error, and the Fatalf method works like Errorf. The difference is that the test function exits immediately after the test failure message is generated. Note that this doesn’t exit all tests; any remaining test functions will execute after the current test function exits.

When should you use Fatal/Fatalf and when should you use Error/Errorf? If the failure of a check in a test means that further checks in the same test function will always fail or cause the test to panic, use Fatal or Fatalf. If you are testing several independent items (such as validating fields in a struct), then use Error or Errorf so you can report as many problems at once. This makes it easier to fix multiple problems without rerunning your tests over and over.

Setting Up and Tearing Down
Sometimes you have some common state that you want to set up before any tests run and remove when testing is complete. Use a TestMain function to manage this state and run your tests:

var testTime time.Time

func TestMain(m *testing.M) {
    fmt.Println("Set up stuff for tests here")
    testTime = time.Now()
    exitVal := m.Run()
    fmt.Println("Clean up stuff after tests here")
    os.Exit(exitVal)
}

func TestFirst(t *testing.T) {
    fmt.Println("TestFirst uses stuff set up in TestMain", testTime)
}

func TestSecond(t *testing.T) {
    fmt.Println("TestSecond also uses stuff set up in TestMain", testTime)
}
Both TestFirst and TestSecond refer to the package-level variable testTime. We declare a function called TestMain with a parameter of type *testing.M. Running go test on a package with a TestMain function calls the function instead of invoking the tests directly. Once the state is configured, call the Run method on *testing.M to run the test functions. The Run method returns the exit code; 0 indicates that all tests passed. Finally, you must call os.Exit with the exit code returned from Run.

Running go test on this produces the output:

$ go test
Set up stuff for tests here
TestFirst uses stuff set up in TestMain 2020-09-01 21:42:36.231508 -0400 EDT
    m=+0.000244286
TestSecond also uses stuff set up in TestMain 2020-09-01 21:42:36.231508 -0400
    EDT m=+0.000244286
PASS
Clean up stuff after tests here
ok      test_examples/testmain  0.006s
Be aware that TestMain is invoked once, not before and after each individual test. Also be aware that you can have only one TestMain per package.

There are two common situations where TestMain is useful:

When you need to set up data in an external repository, such as a database

When the code being tested depends on package-level variables that need to be initialized

As mentioned before (and will be mentioned again!) you should avoid package-level variables in your programs. They make it hard to understand how data flows through your program. If you are using TestMain for this reason, consider refactoring your code.

The Cleanup method on *testing.T is used to clean up temporary resources created for a single test. This method has a single parameter, a function with no input parameters or return values. The function runs when the test completes. For simple tests, you can achieve the same result by using a defer statement, but Cleanup is useful when tests rely on helper functions to set up sample data, like we see in Example 13-1. It’s fine to call Cleanup multiple times. Just like defer, the functions are invoked in last added, first called order.

Example 13-1. Using t.Cleanup
// createFile is a helper function called from multiple tests
func createFile(t *testing.T) (string, error) {
    f, err := os.Create("tempFile")
    if err != nil {
        return "", err
    }
    // write some data to f
    t.Cleanup(func() {
        os.Remove(f.Name())
    })
    return f.Name(), nil
}

func TestFileProcessing(t *testing.T) {
    fName, err := createFile(t)
    if err != nil {
        t.Fatal(err)
    }
    // do testing, don't worry about cleanup
}
Storing Sample Test Data
As go test walks your source code tree, it uses the current package directory as the current working directory. If you want to use sample data to test functions in a package, create a subdirectory named testdata to hold your files. Go reserves this directory name as a place to hold test files. When reading from testdata, always use a relative file reference. Since go test changes the current working directory to the current package, each package accesses its own testdata via a relative file path.

TIP
The text package demonstrates how to use testdata.

Caching Test Results
Just as we learned in Chapter 9 that Go caches compiled packages if they haven’t changed, Go also caches test results when running tests across multiple packages if they have passed and their code hasn’t changed. The tests are recompiled and rerun if you change any file in the package or in the testdata directory. You can also force tests to always run if you pass the flag -count=1 to go test.

Testing Your Public API
The tests that we’ve written are in the same package as the production code. This allows us to test both exported and unexported functions.

If you want to test just the public API of your package, Go has a convention for specifying this. You still keep your test source code in the same directory as the production source code, but you use packagename_test for the package name. Let’s redo our initial test case, using an exported function instead. If we have the following function in the adder package:

func AddNumbers(x, y int) int {
    return x + y
}
then we can test it as public API using the following code in a file in the adder package named adder_public_test.go:

package adder_test

import (
    "testing"
    "test_examples/adder"
)

func TestAddNumbers(t *testing.T) {
    result := adder.AddNumbers(2, 3)
    if result != 5 {
        t.Error("incorrect result: expected 5, got", result)
    }
}
Notice that the package name for our test file is adder_test. We have to import test_examples/adder even though the files are in the same directory. To follow the convention for naming tests, the test function name matches the name of the AddNumbers function. Also note that we use adder.AddNumbers, since we are calling an exported function in a different package.

Just as you can call exported functions from within a package, you can test your public API from a test that is in the same package as your source code. The advantage of using the _test package suffix is that it lets you treat your package as a “black box”; you are forced to interact with it only via its exported functions, methods, types, constants, and variables. Also be aware that you can have test source files with both package names intermixed in the same source directory.

Use go-cmp to Compare Test Results
It can be verbose to write a thorough comparison between two instances of a compound type. While you can use reflect.DeepEqual to compare structs, maps, and slices, there’s a better way. Google released a third-party module called go-cmp that does the comparison for you and returns a detailed description of what does not match. Let’s see how it works by defining a simple struct and a factory function that populates it:

type Person struct {
    Name      string
    Age       int
    DateAdded time.Time
}

func CreatePerson(name string, age int) Person {
    return Person{
        Name:      name,
        Age:       age,
        DateAdded: time.Now(),
    }
}
In our test file, we need to import github.com/google/go-cmp/cmp, and our test function looks like this:

func TestCreatePerson(t *testing.T) {
    expected := Person{
        Name: "Dennis",
        Age:  37,
    }
    result := CreatePerson("Dennis", 37)
    if diff := cmp.Diff(expected, result); diff != "" {
        t.Error(diff)
    }
}
The cmp.Diff function takes in the expected output and the output that was returned by the function that we’re testing. It returns a string that describes any mismatches between the two inputs. If the inputs match, it returns an empty string. We assign the output of the cmp.Diff function to a variable called diff and then check to see if diff is an empty string. If it is not, an error occurred.

We’ll build and run our test and see the output that go-cmp generates when two structs don’t match:

$ go test
--- FAIL: TestCreatePerson (0.00s)
    ch13_cmp_test.go:16:   ch13_cmp.Person{
              Name:      "Dennis",
              Age:       37,
        -     DateAdded: s"0001-01-01 00:00:00 +0000 UTC",
        +     DateAdded: s"2020-03-01 22:53:58.087229 -0500 EST m=+0.001242842",
          }

FAIL
FAIL    ch13_cmp    0.006s
The lines with a - and + indicate the fields whose values differ. Our test failed because our dates didn’t match. This is a problem because we can’t control what date is assigned by the CreatePerson function. We have to ignore the DateAdded field. You do that by specifying a comparator function. Declare the function as a local variable in your test:

comparer := cmp.Comparer(func(x, y Person) bool {
    return x.Name == y.Name && x.Age == y.Age
})
Pass a function to the cmp.Comparer function to create a customer comparator. The function that’s passed in must have two parameters of the same type and return a bool. It also must be symmetric (the order of the parameters doesn’t matter), deterministic (it always returns the same value for the same inputs), and pure (it must not modify its parameters). In our implementation, we are comparing the Name and Age fields and ignoring the DateAdded field.

Then change your call to cmp.Diff to include comparer:

if diff := cmp.Diff(expected, result, comparer); diff != "" {
    t.Error(diff)
}
This is only a quick preview of the most useful features in go-cmp. Check its documentation to learn more about how to control what is compared and the output format.

Table Tests
Most of the time, it takes more than a single test case to validate that a function is working correctly. You could write multiple test functions to validate your function or multiple tests within the same function, but you’ll find that a great deal of the testing logic is repetitive. You set up supporting data and functions, specify inputs, check the outputs, and compare to see if they match your expectations. Rather than writing this over and over, you can take advantage of a pattern called table tests. Let’s take a look at a sample. Assume we have the following function in the table package:

func DoMath(num1, num2 int, op string) (int, error) {
    switch op {
    case "+":
        return num1 + num2, nil
    case "-":
        return num1 - num2, nil
    case "*":
        return num1 + num2, nil
    case "/":
        if num2 == 0 {
            return 0, errors.New("division by zero")
        }
        return num1 / num2, nil
    default:
        return 0, fmt.Errorf("unknown operator %s", op)
    }
}
To test this function, we need to check the different branches, trying out inputs that return valid results, as well as inputs that trigger errors. We could write code like this, but it’s very repetitive:

func TestDoMath(t *testing.T) {
    result, err := DoMath(2, 2, "+")
    if result != 4 {
        t.Error("Should have been 4, got", result)
    }
    if err != nil {
        t.Error("Should have been nil error, got", err)
    }
    result2, err2 := DoMath(2, 2, "-")
    if result2 != 0 {
        t.Error("Should have been 0, got", result2)
    }
    if err2 != nil {
        t.Error("Should have been nil error, got", err2)
    }
    // and so on...
}
Let’s replace this repetition with a table test. First, we declare a slice of anonymous structs. The struct contains fields for the name of the test, the input parameters, and the return values. Each entry in the slice represents another test:

data := []struct {
    name     string
    num1     int
    num2     int
    op       string
    expected int
    errMsg   string
}{
    {"addition", 2, 2, "+", 4, ""},
    {"subtraction", 2, 2, "-", 0, ""},
    {"multiplication", 2, 2, "*", 4, ""},
    {"division", 2, 2, "/", 1, ""},
    {"bad_division", 2, 0, "/", 0, `division by zero`},
}
Next, we loop over each test case in data, invoking the Run method each time. This is the line that does the magic. We pass two parameters to Run, a name for the subtest and a function with a single parameter of type *testing.T. Inside the function, we call DoMath using the fields of the current entry in data, using the same logic over and over. When you run these tests, you’ll see that not only do they pass, but when you use the -v flag, each subtest also now has a name:

for _, d := range data {
    t.Run(d.name, func(t *testing.T) {
        result, err := DoMath(d.num1, d.num2, d.op)
        if result != d.expected {
            t.Errorf("Expected %d, got %d", d.expected, result)
        }
        var errMsg string
        if err != nil {
            errMsg = err.Error()
        }
        if errMsg != d.errMsg {
            t.Errorf("Expected error message `%s`, got `%s`",
                d.errMsg, errMsg)
        }
    })
}
TIP
Comparing error messages can be fragile, because there may not be any compatibility guarantees on the message text. The function that we are testing uses errors.New and fmt.Errorf to make errors, so the only option is to compare the messages. If an error has a custom type, use errors.Is or errors.As to check that the correct error is returned.

Now that we have a way to run lots of tests, let’s learn about code coverage to find out what our tests are testing.

Checking Your Code Coverage
Code coverage is a very useful tool for knowing if you’ve missed any obvious cases. However, reaching 100% code coverage doesn’t guarantee that there aren’t bugs in your code for some inputs. First we’ll see how go test displays code coverage and then we’ll look at the limitations of relying on code coverage alone.

Adding the -cover flag to the go test command calculates coverage information and includes a summary in the test output. If you include a second flag -coverprofile, you can save the coverage information to a file:

go test -v -cover -coverprofile=c.out
If we run our table test with code coverage, the test output now includes a line that indicates the amount of test code coverage, 87.5%. That’s good to know, but it’d be more useful if we could see what we missed. The cover tool included with Go generates an HTML representation of your source code with that information:

go tool cover -html=c.out
When you run it, your web browser should open and show you a page that looks like Figure 13-1.

Initial Code Coverage
Figure 13-1. Initial code coverage
Every file that’s tested appears in the combo box in the upper left. The source code is in one of three colors. Gray is used for lines of code that aren’t testable, green is used for code that’s been covered by a test, and red is used for code that hasn’t been tested. (The reliance on color is unfortunate for readers of the print edition and those who have red-green color blindness. If you are unable to see the colors, the lighter gray is the covered lines.) From looking at this, we can see that we didn’t write a test to cover the default case, when a bad operator is passed to our function. Let’s add that case to our slice of test cases:

{"bad_op", 2, 2, "?", 0, `unknown operator ?`},
When we rerun go test -v -cover -coverprofile=c.out and go tool cover -html=c.out, we see in Figure 13-2 that the final line is covered and we have 100% test code coverage.

Final Code Coverage
Figure 13-2. Final code coverage
Code coverage is a great thing, but it’s not enough. There’s actually a bug in our code, even though we have 100% coverage. Have you noticed it? If not, let’s add another test case and rerun our tests:

{"another_mult", 2, 3, "*", 6, ""},
You should see the error:

table_test.go:57: Expected 6, got 5
There’s a typo in our case for multiplication. It adds the numbers together instead of multiplying them. (Beware the dangers of copy and paste coding!) Fix the code, rerun go test -v -cover -coverprofile=c.out and go tool cover -html=c.out, and you’ll see that tests pass again.

WARNING
Code coverage is necessary, but it’s not sufficient. You can have 100% code coverage and still have bugs in your code!

Benchmarks
Determining how fast (or slow) code runs is surprisingly difficult. Rather than trying to figure it out yourself, you should use the benchmarking support that’s built into Go’s testing framework. Let’s explore it with a function in the test_examples/bench package:

func FileLen(f string, bufsize int) (int, error) {
    file, err := os.Open(f)
    if err != nil {
        return 0, err
    }
    defer file.Close()
    count := 0
    for {
        buf := make([]byte, bufsize)
        num, err := file.Read(buf)
        count += num
        if err != nil {
            break
        }
    }
    return count, nil
}
This function counts the number of characters in a file. It takes in two parameters, the name of the file and the size of the buffer that we are using to read the file (we’ll see the reason for the second parameter in a moment).

Before we see how fast it is, we should test our library to make sure that it works (it does). Here’s a simple test:

func TestFileLen(t *testing.T) {
    result, err := FileLen("testdata/data.txt", 1)
    if err != nil {
        t.Fatal(err)
    }
    if result != 65204 {
        t.Error("Expected 65204, got", result)
    }
}
Now we can see how long it takes our file length function to run. Our goal is to find out what size buffer we should use to read from the file.

NOTE
Before you spend time going down an optimization rabbit hole, be sure that you need to optimize. If your program is already fast enough to meet your responsiveness requirements and is using an acceptable amount of memory, then your time is better spent on adding features and fixing bugs. Your business requirements determine what “fast enough” and “acceptable amount of memory” mean.

In Go, benchmarks are functions in your test files that start with the word Benchmark and take in a single parameter of type *testing.B. This type includes all of the functionality of a *testing.T as well as additional support for benchmarking. Let’s start by looking at a benchmark that uses a buffer size of 1 byte:

var blackhole int

func BenchmarkFileLen1(b *testing.B) {
    for i := 0; i < b.N; i++ {
        result, err := FileLen("testdata/data.txt", 1)
        if err != nil {
            b.Fatal(err)
        }
        blackhole = result
    }
}
The blackhole package-level variable is interesting. We write the results from FileLen to this package-level variable to make sure that the compiler doesn’t get too clever and decide to optimize away the call to FileLen, ruining our benchmark.

Every Go benchmark must have a loop that iterates from 0 to b.N. The testing framework calls our benchmark functions over and over with larger and larger values for N until it is sure that the timing results are accurate. We’ll see this in the output in a moment.

We run a benchmark by passing the -bench flag to go test. This flag expects a regular expression to describe the name of the benchmarks to run. Use -bench=. to run all benchmarks. A second flag, -benchmem, includes memory allocation information in the benchmark output. All tests are run before the benchmarks, so you can only benchmark code when tests pass.

Here’s the output for the benchmark on my computer:

BenchmarkFileLen1-12  25  47201025 ns/op  65342 B/op  65208 allocs/op
Running a benchmark with memory allocation information produces output with five columns. Here’s what each one means:

BenchmarkFileLen1-12
The name of the benchmark, a hyphen, and the value of GOMAXPROCS for the benchmark.

25
The number of times that the test ran to produce a stable result.

47201025 ns/op
How long it took to run a single pass of this benchmark, in nanoseconds (there are 1,000,000,000 nanoseconds in a second).

65342 B/op
The number of bytes allocated during a single pass of the benchmark.

65208 allocs/op
The number of times bytes had to be allocated from the heap during a single pass of the benchmark. This will always be less than or equal to the number of bytes allocated.

Now that we have results for a buffer of 1 byte, let’s see what the results look like when we use buffers of different sizes:

func BenchmarkFileLen(b *testing.B) {
    for _, v := range []int{1, 10, 100, 1000, 10000, 100000} {
        b.Run(fmt.Sprintf("FileLen-%d", v), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                result, err := FileLen("testdata/data.txt", v)
                if err != nil {
                    b.Fatal(err)
                }
                blackhole = result
            }
        })
    }
}
Just like we launched table tests using t.Run, we’re using b.Run to launch benchmarks that only vary based on input. Here are the results of this benchmark on my computer:

BenchmarkFileLen/FileLen-1-12          25  47828842 ns/op   65342 B/op  65208 allocs/op
BenchmarkFileLen/FileLen-10-12        230   5136839 ns/op  104488 B/op   6525 allocs/op
BenchmarkFileLen/FileLen-100-12      2246    509619 ns/op   73384 B/op    657 allocs/op
BenchmarkFileLen/FileLen-1000-12    16491     71281 ns/op   68744 B/op     70 allocs/op
BenchmarkFileLen/FileLen-10000-12   42468     26600 ns/op   82056 B/op     11 allocs/op
BenchmarkFileLen/FileLen-100000-12  36700     30473 ns/op  213128 B/op      5 allocs/op
These results aren’t surprising; as we increase the size of the buffer, we make fewer allocations and our code runs faster, until the buffer is bigger than the file. When the buffer is bigger than the size of the file, there are extra allocations that slow down the output. If we expect files of roughly this size, a buffer of 10,000 bytes would work best.

But there’s a change we can make that improves the numbers more. We are reallocating the buffer every time we get the next set of bytes from the file. That’s unnecessary. If we move the byte slice allocation before the loop and rerun our benchmark, we see an improvement:

BenchmarkFileLen/FileLen-1-12          25  46167597 ns/op     137 B/op  4 allocs/op
BenchmarkFileLen/FileLen-10-12        261   4592019 ns/op     152 B/op  4 allocs/op
BenchmarkFileLen/FileLen-100-12      2518    478838 ns/op     248 B/op  4 allocs/op
BenchmarkFileLen/FileLen-1000-12    20059     60150 ns/op    1160 B/op  4 allocs/op
BenchmarkFileLen/FileLen-10000-12   62992     19000 ns/op   10376 B/op  4 allocs/op
BenchmarkFileLen/FileLen-100000-12  51928     21275 ns/op  106632 B/op  4 allocs/op
The number of allocations are now consistent and small, just four allocations for every buffer size. What is interesting is that we now can make trade-offs. If we are tight on memory, we can use a smaller buffer size and save memory at the expense of performance.

PROFILING YOUR GO CODE
If benchmarking reveals that you have a performance or memory problem, the next step is figuring out exactly what the problem is. Go includes profiling tools that gather CPU and memory usage data from a running program as well as tools that help you visualize and interpret the generated data. You can even expose a web service endpoint to gather profiling information remotely from a running Go service.

Discussing the profiler is a topic that’s beyond the scope of this book. There are many great resources available online with information on it. A good starting point is the blog post Profiling Go programs with pprof by Julia Evans.

Stubs in Go
So far, we’ve written tests for functions that didn’t depend on other code. This is not typical as most code is filled with dependencies. As we saw in Chapter 7, there are two ways that Go allows us to abstract function calls: defining a function type and defining an interface. These abstractions not only help us write modular production code; they also help us write unit tests.

TIP
When your code depends on abstractions, it’s easier to write unit tests!

Lets take a look at an example in the test_examples/solver package. We define a type called Processor:

type Processor struct {
    Solver MathSolver
}
It has a field of type MathSolver:

type MathSolver interface {
    Resolve(ctx context.Context, expression string) (float64, error)
}
We’ll implement and test MathSolver in a bit.

Processor also has a method that reads an expression from an io.Reader and returns the calculated value:

func (p Processor) ProcessExpression(ctx context.Context, r io.Reader)
                                    (float64, error) {
    curExpression, err := readToNewLine(r)
    if err != nil {
        return 0, err
    }
    if len(curExpression) == 0 {
        return 0, errors.New("no expression to read")
    }
    answer, err := p.Solver.Resolve(ctx, curExpression)
    return answer, err
}
Let’s write the code to test ProcessExpression. First, we need a simple implementation of the Resolve method to write our test:

type MathSolverStub struct{}

func (ms MathSolverStub) Resolve(ctx context.Context, expr string)
                                (float64, error) {
    switch expr {
    case "2 + 2 * 10":
        return 22, nil
    case "( 2 + 2 ) * 10":
        return 40, nil
    case "( 2 + 2 * 10":
        return 0, errors.New("invalid expression: ( 2 + 2 * 10")
    }
    return 0, nil
}
Next, we write a unit test that uses this stub (production code should test the error messages too, but for the sake of brevity, we’ll leave those out):

func TestProcessorProcessExpression(t *testing.T) {
    p := Processor{MathSolverStub{}}
    in := strings.NewReader(`2 + 2 * 10
( 2 + 2 ) * 10
( 2 + 2 * 10`)
    data := []float64{22, 40, 0}
    hasErr := []bool{false, false, true}
    for i, d := range data {
        result, err := p.ProcessExpression(context.Background(), in)
        if err != nil && !hasErr[i] {
            t.Error(err)
        }
        if result != d {
            t.Errorf("Expected result %f, got %f", d, result)
        }
    }
}
We can then run our test and see that everything works.

While most Go interfaces only specify one or two methods, this isn’t always the case. You sometimes find yourself with an interface that has many methods. Assume you have an interface that looks like this:

type Entities interface {
    GetUser(id string) (User, error)
    GetPets(userID string) ([]Pet, error)
    GetChildren(userID string) ([]Person, error)
    GetFriends(userID string) ([]Person, error)
    SaveUser(user User) error
}
There are two patterns for testing code that depends on large interfaces. The first is to embed the interface in a struct. Embedding an interface in a struct automatically defines all of the interface’s methods on the struct. It doesn’t provide any implementations of those methods, so you need to implement the methods that you care about for the current test. Let’s assume that Logic is a struct that has a field of type Entities:

type Logic struct {
    Entities Entities
}
Assume you want to test this method:

func (l Logic) GetPetNames(userId string) ([]string, error) {
    pets, err := l.Entities.GetPets(userId)
    if err != nil {
        return nil, err
    }
    out := make([]string, len(pets))
    for _, p := range pets {
        out = append(out, p.Name)
    }
    return out, nil
}
This method uses only one of the methods declared on Entities, GetPets. Rather than creating a stub that implements every single method on Entities just to test GetPets, you can write a stub struct that only implements the method you need to test this method:

type GetPetNamesStub struct {
    Entities
}

func (ps GetPetNamesStub) GetPets(userID string) ([]Pet, error) {
    switch userID {
    case "1":
        return []Pet{{Name: "Bubbles"}}, nil
    case "2":
        return []Pet{{Name: "Stampy"}, {Name: "Snowball II"}}, nil
    default:
        return nil, fmt.Errorf("invalid id: %s", userID)
    }
}
We then write our unit test, with our stub injected into Logic:

func TestLogicGetPetNames(t *testing.T) {
    data := []struct {
        name     string
        userID   string
        petNames []string
    }{
        {"case1", "1", []string{"Bubbles"}},
        {"case2", "2", []string{"Stampy", "Snowball II"}},
        {"case3", "3", nil},
    }
    l := Logic{GetPetNamesStub{}}
    for _, d := range data {
        t.Run(d.name, func(t *testing.T) {
            petNames, err := l.GetPetNames(d.userID)
            if err != nil {
                t.Error(err)
            }
            if diff := cmp.Diff(d.petNames, petNames); diff != "" {
                t.Error(diff)
            }
        })
    }
}
(By the way, the GetPetNames method has a bug. Did you see it? Even simple methods can sometimes have bugs.)

WARNING
If you embed an interface in a stub struct, make sure you provide an implementation for every method that’s called during your test! If you call an unimplemented method, your tests will panic.

If you need to implement only one or two methods in an interface for a single test, this technique works well. The drawback comes when you need to call the same method in different tests with different inputs and outputs. When that happens, you need to either include every possible result for every test within the same implementation or reimplement the struct for each test. This quickly becomes difficult to understand and maintain. A better solution is to create a stub struct that proxies method calls to function fields. For each method defined on Entities, we define a function field with a matching signature on our stub struct:

type EntitiesStub struct {
    getUser     func(id string) (User, error)
    getPets     func(userID string) ([]Pet, error)
    getChildren func(userID string) ([]Person, error)
    getFriends  func(userID string) ([]Person, error)
    saveUser    func(user User) error
}
We then make EntitiesStub meet the Entities interface by defining the methods. In each method, we invoke the associated function field. For example:

func (es EntitiesStub) GetUser(id string) (User, error) {
    return es.getUser(id)
}

func (es EntitiesStub) GetPets(userID string) ([]Pet, error) {
    return es.getPets(userID)
}
Once you create this stub, you can supply different implementations of different methods in different test cases via the fields in the data struct for a table test:

func TestLogicGetPetNames(t *testing.T) {
    data := []struct {
        name     string
        getPets  func(userID string) ([]Pet, error)
        userID   string
        petNames []string
        errMsg   string
    }{
        {"case1", func(userID string) ([]Pet, error) {
            return []Pet{{Name: "Bubbles"}}, nil
        }, "1", []string{"Bubbles"}, ""},
        {"case2", func(userID string) ([]Pet, error) {
            return nil, errors.New("invalid id: 3")
        }, "3", nil, "invalid id: 3"},
    }
    l := Logic{}
    for _, d := range data {
        t.Run(d.name, func(t *testing.T) {
            l.Entities = EntitiesStub{getPets: d.getPets}
            petNames, err := l.GetPetNames(d.userID)
            if diff := cmp.Diff(petNames, d.petNames); diff != "" {
                t.Error(diff)
            }
            var errMsg string
            if err != nil {
                errMsg = err.Error()
            }
            if errMsg != d.errMsg {
                t.Errorf("Expected error `%s`, got `%s`", d.errMsg, errMsg)
            }
        })
    }
}
We add a field of function type to data’s anonymous struct. In each test case, we specify a function that returns the data that GetPets would return. When you write your test stubs this way, it’s clear what the stubs should return for each test case. As each test runs, we instantiate a new EntitiesStub and assign the getPets function field in our test data to the getPets function field in EntitiesStub.

MOCKS AND STUBS
The terms mock and stub are often used interchangeably, but they are actually two different concepts. Martin Fowler, a respected voice on all things related to software development, wrote a blog post on mocks that, among other things, covers the difference between mocks and stubs. In short, a stub returns a canned value for a given input, whereas a mock validates that a set of calls happen in the expected order with the expected inputs.

We used stubs in our examples to return canned values to a given response. You can write your own mocks by hand, or you can use a third-party library to generate them. The two most popular options are the gomock library from Google and the testify library from Stretchr, Inc.

httptest
It can be difficult to write tests for a function that calls an HTTP service. Traditionally, this became an integration test, requiring you to stand up a test instance of the service that the function calls. The Go standard library includes the net/http/httptest package to make it easier to stub HTTP services. Let’s go back to our test_examples/solver package and provide an implementation of MathSolver that calls an HTTP service to evaluate expressions:

type RemoteSolver struct {
    MathServerURL string
    Client        *http.Client
}

func (rs RemoteSolver) Resolve(ctx context.Context, expression string)
                              (float64, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet,
        rs.MathServerURL+"?expression="+url.QueryEscape(expression),
        nil)
    if err != nil {
        return 0, err
    }
    resp, err := rs.Client.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    contents, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return 0, err
    }
    if resp.StatusCode != http.StatusOK {
        return 0, errors.New(string(contents))
    }
    result, err := strconv.ParseFloat(string(contents), 64)
    if err != nil {
        return 0, err
    }
    return result, nil
}
Now let’s see how to use the httptest library to test this code without standing up a server. The code is in the TestRemoteSolver_Resolve in test_examples/solver/remote_solver_test.go, but here are the highlights. First, we want to make sure that the data that’s passed into the function arrives on the server. So in our test function, we define a type called info to hold our input and output and a variable called io that is assigned the current input and output:

type info struct {
    expression string
    code       int
    body       string
}
var io info
Next, we set up our fake remote server and use it to configure an instance of RemoteSolver:

server := httptest.NewServer(
    http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        expression := req.URL.Query().Get("expression")
        if expression != io.expression {
            rw.WriteHeader(http.StatusBadRequest)
            fmt.Fprintf(rw, "expected expression '%s', got '%s'",
                io.expression, expression)
            return
        }
        rw.WriteHeader(io.code)
        rw.Write([]byte(io.body))
    }))
defer server.Close()
rs := RemoteSolver{
    MathServerURL: server.URL,
    Client:        server.Client(),
}
The httptest.NewServer function creates and starts an HTTP server on a random unused port. You need to provide an http.Handler implementation to process the request. Since this is a server, you must close it when the test completes. The server instance has its URL specified in the URL field of the server instance and a preconfigured http.Client for communicating with the test server. We pass these into RemoteSolver.

The rest of the function works like every other table test that we’ve seen:

    data := []struct {
        name   string
        io     info
        result float64
    }{
        {"case1", info{"2 + 2 * 10", http.StatusOK, "22"}, 22},
        // remaining cases
    }
    for _, d := range data {
        t.Run(d.name, func(t *testing.T) {
            io = d.io
            result, err := rs.Resolve(context.Background(), d.io.expression)
            if result != d.result {
                t.Errorf("io `%f`, got `%f`", d.result, result)
            }
            var errMsg string
            if err != nil {
                errMsg = err.Error()
            }
            if errMsg != d.errMsg {
                t.Errorf("io error `%s`, got `%s`", d.errMsg, errMsg)
            }
        })
    }
The interesting thing to note is that the variable io has been captured by two different closures: the one for the stub server and the one for running each test. We write to it in one closure and read it in the other. This is a bad idea in production code, but it works well in test code within a single function.

Integration Tests and Build Tags
Even though httptest provides a way to avoid testing against external services, you should still write integration tests, automated tests that connect to other services. These validate that your understanding of the service’s API is correct. The challenge is figuring out how to group your automated tests; you only want to run integration tests when the support environment is present. Also, integration tests tend to be slower than unit tests, so they are usually run less frequently.

The Go compiler provides build tags to control when code is compiled. Build tags are specified on the first line of a file with a magic comment that starts with // +build. The original intent for build tags was to allow different code to be compiled on different platforms, but they are also useful for splitting tests into groups. Tests in files without build tags run all the time. These are the unit tests that don’t have dependencies on external resources. Tests in files with a build tag are only run when the supporting resources are available.

Let’s try this out with our math solving project. Use Docker to download a server implementation with docker pull jonbodner/math-server and then run the server locally on port 8080 with docker run -p 8080:8080 jonbodner/math-server.

NOTE
If you don’t have Docker installed or if you want to build the code for yourself, you can find it on GitHub.

We  need  to  write  an  integration  test  to  make  sure  that  our  Resolve method properly  communicates with the math server. The test_examples/solver/remote_solver_integration_test.go file has a complete test in the TestRemoteSolver_ResolveIntegration function. The test looks like every other table test that we’ve written. The interesting thing is the first line of the file, separated from the package declaration by a newline is:

// +build integration
To run our integration test alongside the other tests we’ve written, use:

$ go test -tags integration -v ./...
USING THE -SHORT FLAG
Another option is to use go test with the -short flag. If you want to skip over tests that take a long time, label your slow tests by placing the the following code at the start of the test function:

    if testing.Short() {
        t.Skip("skipping test in short mode.")
    }
When you want to run only short tests, pass the -short flag to go test.

There are a few problems with the -short flag. If you use it, there are only two levels of testing: short tests and all tests. By using build tags, you can group your integration tests, specifying which service they need in order to run. Another argument against using the -short flag to indicate integration tests is philosophical. Build tags indicate a dependency, while the -short flag is only meant to indicate that you don’t want to run tests that take a long time. Those are different concepts. Finally, I find the -short flag unintuitive. You should run short tests all the time. It makes more sense to require a flag to include long-running tests, not to exclude them.

Finding Concurrency Problems with the Race Checker
Even with Go’s built-in support for concurrency, bugs still happen. It’s easy to accidentally reference a variable from two different goroutines without acquiring a lock. The computer science term for this is a data race. To help find these sorts of bugs, Go includes a race checker. It isn’t guaranteed to find every single data race in your code, but if it finds one, you should put proper locks around what it finds.

Let’s look at a simple example in test_examples/race/race.go:

func getCounter() int {
    var counter int
    var wg sync.WaitGroup
    wg.Add(5)
    for i := 0; i < 5; i++ {
        go func() {
            for i := 0; i < 1000; i++ {
                counter++
            }
            wg.Done()
        }()
    }
    wg.Wait()
    return counter
}
This code launches five goroutines, has each of them update a shared counter variable 1000 times, and then returns the result. You’d expect it to be 5000, so let’s verify this with a unit test in test_examples/race/race_test.go:

func TestGetCounter(t *testing.T) {
    counter := getCounter()
    if counter != 5000 {
        t.Error("unexpected counter:", counter)
    }
}
If you run go test a few times, you’ll see that sometimes it passes, but most of the time it fails with an error message like:

unexpected counter: 3673
The problem is that there’s a data race in the code. In a program this simple, the cause is obvious: multiple goroutines are trying to update counter simultaneously and some of their updates are lost. In more complicated programs, these sorts of races are harder to see. Let’s see what the race checker does. Use the flag -race with go test to enable it:

$ go test -race
==================
WARNING: DATA RACE
Read at 0x00c000128070 by goroutine 10:
  test_examples/race.getCounter.func1()
      test_examples/race/race.go:12 +0x45

Previous write at 0x00c000128070 by goroutine 8:
  test_examples/race.getCounter.func1()
      test_examples/race/race.go:12 +0x5b
The traces make it clear that the line counter++ is the source of our problems.

WARNING
Some people try to fix race conditions by inserting “sleeps” into their code, trying to space out access to the variable that’s being accessed by multiple goroutines. This is a bad idea. Doing so might eliminate the problem in some cases, but the code is still wrong and it will fail in some situations.

You can also use the -race flag when you build your programs. This creates a binary that includes the race checker and that reports any races it finds to the console. This allows you to find data races in code that doesn’t have tests.

If the race checker is so useful, why isn’t it enabled all the time for testing and production? A binary with -race enabled runs approximately ten times slower than a normal binary. That isn’t a problem for test suites that take a second to run, but for large test suites that take several minutes, a 10x slowdown reduces productivity.

Wrapping Up
In this chapter, we’ve learned how to write tests and improve code quality using Go’s built-in support for testing, code coverage, benchmarking, and data race checking. In the next chapter, we’re going to explore some Go features that allow you to break the rules: the unsafe package, reflection, and cgo.