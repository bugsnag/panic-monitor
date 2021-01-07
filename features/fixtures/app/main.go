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
	}
}
