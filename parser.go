package main

func panicHeaders() [][]byte {
	return [][]byte{
		[]byte("panic:"),
		[]byte("fatal error:"),
	}
}
