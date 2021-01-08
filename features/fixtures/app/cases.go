package main

import (
	"fmt"
	"os"
)

func concurrentReadWrite() {
	m := map[string]int{}

	go func() {
		for {
			m["x"] = 1
		}
	}()
	for {
		_ = m["x"]
	}
}

func explicitPanic() {
	panic("PANIQ!")
}

func nilGoroutine() {
	var f func()
	go f()
}

func fakePanicRealPanic() {
	fmt.Fprint(os.Stderr, "panic: foo\n\n")
	for i := 0; i < 1024; i++ {
		fmt.Fprint(os.Stderr, "foobarbaz")
	}
	os.Stderr.Sync()

	panic("REAL PANIC!")
}
