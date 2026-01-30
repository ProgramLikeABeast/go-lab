package main

import (
	"fmt"
	"os"
	"runtime/debug"
)

// =============================================================================
// DEFER BASICS: Schedule cleanup to run when function returns
// =============================================================================

func demoBasicDefer() {
	fmt.Println("\n=== BASIC DEFER ===")

	fmt.Println("First")
	defer fmt.Println("Deferred (runs at function end)")
	fmt.Println("Second")
	fmt.Println("Third")
	// Output order: First, Second, Third, Deferred
}

// =============================================================================
// DEFER ORDER: LIFO (Last In, First Out) - like a stack
// =============================================================================

func demoDeferOrder() {
	fmt.Println("\n=== DEFER ORDER (LIFO) ===")

	defer fmt.Println("1st defer")
	defer fmt.Println("2nd defer")
	defer fmt.Println("3rd defer")
	fmt.Println("Function body")
	// Output: Function body, 3rd defer, 2nd defer, 1st defer
}

// =============================================================================
// DEFER WITH LOOPS: Common gotcha!
// =============================================================================

func demoDeferInLoop() {
	fmt.Println("\n=== DEFER IN LOOPS (GOTCHA) ===")

	// WRONG: All defers run at function end, not loop iteration end
	fmt.Println("Wrong way (don't do this with file handles):")
	for i := 0; i < 3; i++ {
		defer fmt.Printf("  Loop defer %d\n", i)
	}

	// RIGHT: Use closure or separate function for loop cleanup
	fmt.Println("\nRight way (wrap in function):")
	for i := 0; i < 3; i++ {
		func(n int) {
			defer fmt.Printf("  Closure defer %d\n", n)
			// Deferred statement runs when this closure returns
		}(i)
	}
}

// =============================================================================
// DEFER ARGUMENT EVALUATION: Arguments captured at defer time!
// =============================================================================

func demoDeferArgs() {
	fmt.Println("\n=== DEFER ARGUMENT EVALUATION ===")

	x := 10
	defer fmt.Printf("Deferred with x=%d (captured at defer time)\n", x)
	x = 20
	fmt.Printf("Current x=%d\n", x)
	// Output: Current x=20, then Deferred with x=10

	// To capture later value, use closure
	y := 10
	defer func() {
		fmt.Printf("Deferred closure with y=%d (evaluated at call time)\n", y)
	}()
	y = 20
	// Output: y=20 (closure captures variable, not value)
}

// =============================================================================
// DEFER WITH NAMED RETURNS: Can modify return values!
// =============================================================================

func demoNamedReturns() (result int) {
	defer func() {
		result += 10 // Modifies the named return value
	}()
	return 5 // Actually returns 15!
}

func demoDeferNamedReturns() {
	fmt.Println("\n=== DEFER WITH NAMED RETURNS ===")
	fmt.Println("demoNamedReturns() =", demoNamedReturns())
}

// =============================================================================
// PRACTICAL DEFER: Resource cleanup
// =============================================================================

func demoResourceCleanup() {
	fmt.Println("\n=== PRACTICAL DEFER: RESOURCE CLEANUP ===")

	// File handling pattern
	fmt.Println("File cleanup pattern:")
	fmt.Println(`
    func processFile(path string) error {
        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()  // Guaranteed to run, even on error
        
        // Process file...
        return nil
    }`)

	// Mutex pattern
	fmt.Println("\nMutex cleanup pattern:")
	fmt.Println(`
    func (c *Counter) Increment() {
        c.mu.Lock()
        defer c.mu.Unlock()  // Always unlocks, even if panic
        c.count++
    }`)

	// Timing pattern
	fmt.Println("\nTiming pattern:")
	fmt.Println(`
    func slowOperation() {
        defer trackTime(time.Now())  // Prints duration at end
        // ... operation ...
    }
    
    func trackTime(start time.Time) {
        fmt.Printf("Took %v\n", time.Since(start))
    }`)
}

// =============================================================================
// PANIC: Abnormal termination (like throw in other languages)
// =============================================================================

func causePanic() {
	fmt.Println("About to panic...")
	panic("something went wrong!")
	fmt.Println("This never prints")
}

func demoPanic() {
	fmt.Println("\n=== PANIC ===")
	fmt.Println("Panic causes immediate function termination")
	fmt.Println("Stack unwinds, running defers along the way")
	fmt.Println("If not recovered, program crashes with stack trace")
	fmt.Println("\n(We'll demonstrate with recover below)")
}

// =============================================================================
// RECOVER: Catch panics (like catch in other languages)
// =============================================================================

func safeOperation() {
	// Recover only works inside a deferred function
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from: %v\n", r)
		}
	}()

	fmt.Println("Starting risky operation...")
	panic("unexpected error!")
	fmt.Println("This never executes")
}

func demoRecover() {
	fmt.Println("\n=== RECOVER ===")

	safeOperation()
	fmt.Println("Execution continues after recovery!")
}

// =============================================================================
// PANIC + RECOVER PATTERNS
// =============================================================================

// Pattern 1: Convert panic to error (useful in libraries)
func safelyDo(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()
	f()
	return nil
}

// Pattern 2: Log panic with stack trace
func logPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: %v\nStack trace:\n%s\n", r, debug.Stack())
		}
	}()
	panic("something bad!")
}

// Pattern 3: HTTP server pattern (recover per request)
func handleRequest() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("500 Internal Server Error")
			// Log error, don't crash server
		}
	}()
	// Process request...
	panic("database error")
}

func demoPanicPatterns() {
	fmt.Println("\n=== PANIC/RECOVER PATTERNS ===")

	// Pattern 1
	fmt.Println("Pattern 1: Convert panic to error")
	err := safelyDo(func() {
		panic("oops!")
	})
	fmt.Println("Error:", err)

	// Pattern 2 (uncomment to see stack trace)
	// fmt.Println("\nPattern 2: Log with stack trace")
	// logPanic()

	// Pattern 3
	fmt.Println("\nPattern 3: HTTP server pattern")
	handleRequest()
	fmt.Println("Server still running!")
}

// =============================================================================
// WHEN TO USE PANIC (and when NOT to)
// =============================================================================

func demoWhenToPanic() {
	fmt.Println("\n=== WHEN TO USE PANIC ===")

	fmt.Println(`
Use panic for:
  - Truly unrecoverable errors (corrupted state)
  - Programming errors that should never happen
  - Initialization failures (can't proceed anyway)
  - "Must" helper functions (MustCompile, MustParse)

DON'T panic for:
  - Expected errors (file not found, network timeout)
  - User input validation
  - API errors
  - Anything that should be communicated to caller

Go philosophy: "Don't panic!"
Return errors for recoverable situations.
`)

	// Example: Must pattern for init-time errors
	mustGetEnv := func(key string) string {
		val := os.Getenv(key)
		if val == "" {
			panic(fmt.Sprintf("required env var %s not set", key))
		}
		return val
	}

	// Safe to call at init time
	_ = mustGetEnv
	fmt.Println("Must pattern: acceptable for startup requirements")
}

// =============================================================================
// DEFER + PANIC INTERACTION
// =============================================================================

func demoDeferPanicInteraction() {
	fmt.Println("\n=== DEFER + PANIC INTERACTION ===")

	defer fmt.Println("1. First defer (runs last)")
	defer fmt.Println("2. Second defer")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("3. Recovered:", r)
		}
	}()
	defer fmt.Println("4. Fourth defer (runs before recover)")

	panic("test panic")

	// Order: 4, recover, 2, 1
}

func main() {
	fmt.Println("========================================")
	fmt.Println("GO DEFER, PANIC, AND RECOVER")
	fmt.Println("========================================")

	demoBasicDefer()
	demoDeferOrder()
	demoDeferInLoop()
	demoDeferArgs()
	demoDeferNamedReturns()
	demoResourceCleanup()
	demoPanic()
	demoRecover()
	demoPanicPatterns()
	demoWhenToPanic()
	demoDeferPanicInteraction()
}
