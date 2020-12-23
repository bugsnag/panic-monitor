package main

func main() {
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
