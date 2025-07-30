package main

// extern int raise(int sig);
import "C"

func doSegfault() {
	C.raise(11) // SIGSEGV
}

