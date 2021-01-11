package main

import (
	"os"
)

func main() {
	testcase := os.Getenv("TESTCASE")
	switch testcase {
	case "explicit panic":
		explicitPanic()
	case "concurrent read/write":
		concurrentReadWrite()
	case "nil goroutine":
		nilGoroutine()
	case "stack exhaustion":
		stackExhaustion()
	case "fake panic then real":
		fakePanicRealPanic()
	case "array overflow":
		arrayOverflow()
	case "nil pointer deref":
		nilPointerDeref()
	case "bad reflect swap":
		badSwap(0, "boo")
	case "goroutine":
		panicInAGoroutine()
	}
}
