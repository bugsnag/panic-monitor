package main

import (
	"reflect"
	"time"
)

func badSwap(index int, slice interface{}) {
	reflect.Swapper(slice)(0, 1)
}

func panicInAGoroutine() {
	go func() {
		panic("at the disco?")
	}()
	// give it a sec to get there
	time.Sleep(time.Millisecond * 100)
}
