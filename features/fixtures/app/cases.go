package main

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
