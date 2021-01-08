package main

import (
	"reflect"
)

func badSwap(index int, slice interface{}) {
	reflect.Swapper(slice)(0, 1)
}
