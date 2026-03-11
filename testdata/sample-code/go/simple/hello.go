package main

import (
	"fmt"
	"runtime"
	"time"
)

// Constants
const (
	greeting    = "Hello"
	defaultName = "World"
	version     = "1.0.0"
)

// Variables
var (
	startTime = time.Now()
	counter   = 0
	debug     = true
)

// Function with multiple returns
func getGreeting() (string, string) {
	return greeting, defaultName
}

// Function with named returns
func getVersion() (major, minor, patch int) {
	major = 1
	minor = 0
	patch = 0
	return
}

// Function with parameters
func greet(name string) string {
	counter++
	if name == "" {
		name = defaultName
	}
	return fmt.Sprintf("%s, %s!", greeting, name)
}

// Function with error handling
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Function with variadic parameters
func sum(numbers ...int) int {
	total := 0
	for _, n := range numbers {
		total += n
	}
	return total
}

// Function with defer
func deferExample() {
	defer fmt.Println("This runs last")
	defer fmt.Println("This runs second")
	fmt.Println("This runs first")
}

// Function with panic/recover
func safeDivide(a, b int) (result int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			result = 0
		}
	}()
	return a / b
}

// Function with closure
func counterClosure() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// Function with method receiver
type person struct {
	name string
	age  int
}

func (p person) greet() string {
	return fmt.Sprintf("Hi, I'm %s, %d years old", p.name, p.age)
}

func (p *person) haveBirthday() {
	p.age++
}

// Interface example
type greeter interface {
	greet() string
}

func printGreeting(g greeter) {
	fmt.Println(g.greet())
}

// Function with type switch
func getType(v interface{}) string {
	switch t := v.(type) {
	case int:
		return fmt.Sprintf("int: %d", t)
	case string:
		return fmt.Sprintf("string: %s", t)
	case bool:
		return fmt.Sprintf("bool: %v", t)
	default:
		return fmt.Sprintf("unknown type: %T", v)
	}
}

// Function with select statement
func selectExample() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	case <-time.After(150 * time.Millisecond):
		fmt.Println("timeout")
	}
}

// Function with for loop
func loopExample() {
	// Traditional for loop
	for i := 0; i < 3; i++ {
		fmt.Printf("i = %d\n", i)
	}

	// While-style loop
	j := 0
	for j < 3 {
		fmt.Printf("j = %d\n", j)
		j++
	}

	// Infinite loop with break
	k := 0
	for {
		if k >= 3 {
			break
		}
		fmt.Printf("k = %d\n", k)
		k++
	}

	// Range loop
	nums := []int{1, 2, 3}
	for idx, val := range nums {
		fmt.Printf("nums[%d] = %d\n", idx, val)
	}
}

// Function with slice operations
func sliceExample() {
	// Create slice
	nums := make([]int, 0, 5)

	// Append
	nums = append(nums, 1, 2, 3)

	// Slice operations
	sub := nums[1:3]

	// Copy
	copy := make([]int, len(nums))
	copy(copy, nums)

	fmt.Printf("nums: %v, sub: %v, copy: %v\n", nums, sub, copy)
}

// Function with map operations
func mapExample() {
	// Create map
	ages := make(map[string]int)

	// Add entries
	ages["Alice"] = 30
	ages["Bob"] = 25

	// Check existence
	if age, ok := ages["Alice"]; ok {
		fmt.Printf("Alice is %d years old\n", age)
	}

	// Delete
	delete(ages, "Bob")

	// Iterate
	for name, age := range ages {
		fmt.Printf("%s is %d\n", name, age)
	}
}

// Function with pointer
func pointerExample() {
	x := 42
	p := &x

	fmt.Printf("x = %d, *p = %d\n", x, *p)

	*p = 100
	fmt.Printf("x = %d (after pointer assignment)\n", x)
}

// Main function
func main() {
	// Basic function calls
	fmt.Println(greet("Alice"))
	fmt.Println(greet(""))

	// Multiple returns
	g, n := getGreeting()
	fmt.Printf("%s, %s\n", g, n)

	// Named returns
	maj, min, pat := getVersion()
	fmt.Printf("Version: %d.%d.%d\n", maj, min, pat)

	// Error handling
	if res, err := divide(10, 2); err == nil {
		fmt.Printf("10 / 2 = %f\n", res)
	}
	if _, err := divide(10, 0); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Variadic function
	fmt.Printf("Sum: %d\n", sum(1, 2, 3, 4, 5))

	// Defer example
	deferExample()

	// Panic/recover
	fmt.Printf("Safe divide: %d\n", safeDivide(10, 0))

	// Closure
	counter := counterClosure()
	fmt.Println(counter())
	fmt.Println(counter())
	fmt.Println(counter())

	// Methods
	p := person{name: "Alice", age: 30}
	fmt.Println(p.greet())
	p.haveBirthday()
	fmt.Printf("After birthday: %d\n", p.age)

	// Interface
	printGreeting(p)

	// Type switch
	fmt.Println(getType(42))
	fmt.Println(getType("hello"))
	fmt.Println(getType(true))
	fmt.Println(getType(3.14))

	// Select
	selectExample()

	// Loops
	loopExample()

	// Slices
	sliceExample()

	// Maps
	mapExample()

	// Pointers
	pointerExample()

	// Print runtime info
	fmt.Printf("Go version: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("Program ran for: %v\n", time.Since(startTime))
}
