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

func stackExhaustion() {
	var f func(a [1000]int64)
	f = func(a [1000]int64) {
		f(a)
	}
	f([1000]int64{})
}

func arrayOverflow() {
	items := []int{1,2}
	fmt.Printf("the third item: %d", items[3])
}

func nilPointerDeref() {
	type Item struct {
		Text *string
	}
	item := Item{}
	fmt.Printf("contents: %s", *item.Text)
}
