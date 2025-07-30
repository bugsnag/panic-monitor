package main

import (
	"fmt"
	"os"

	"github.com/bugsnag/bugsnag-go/v2"
)

func main() {
	bugsnag.Configure(bugsnag.Configuration{})
	testcase := os.Getenv("TESTCASE")
	switch testcase {
	case "explicit panic":
		panic("oh no!")
	default:
		fmt.Printf("unknown case: '%s'", testcase)
	}
}
