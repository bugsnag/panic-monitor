package main

import "testing"

func TestFindPanic(t *testing.T) {
	reader := NewOutputReader()
	if reader.FoundPanic() {
		t.Fatalf("erroneously found panic\n")
	}

	reader.ReadStderr([]byte("don't panic: its gonna be ok"))
	if !reader.FoundPanic() {
		t.Fatalf("failed to find expected 'panic: its gonna be ok'\n")
	}
	if reader.panicType < 0 {
		t.Fatalf("failed to set panic type\n")
	}
	panicType := string(panicHeaders[reader.panicType])
	if panicType != "panic:" {
		t.Fatalf("incorrect panic type: '%s'\n", panicType)
	}
}

func TestFindFatalError(t *testing.T) {
	reader := NewOutputReader()
	if reader.FoundPanic() {
		t.Fatalf("erroneously found panic\n")
	}

	reader.ReadStderr([]byte("fatal error: broken pipe"))
	if !reader.FoundPanic() {
		t.Fatalf("failed to find expected 'fatal error: broken pipe'\n")
	}
	if reader.panicType < 0 {
		t.Fatalf("failed to set panic type\n")
	}
	panicType := string(panicHeaders[reader.panicType])
	if panicType != "fatal error:" {
		t.Fatalf("incorrect panic type: '%s'\n", panicType)
	}
}

func TestFindPanicAcrossWrites(t *testing.T) {
	reader := NewOutputReader()
	if reader.FoundPanic() {
		t.Fatalf("erroneously found panic\n")
	}
	reader.ReadStderr([]byte("not\npan"))
	if reader.FoundPanic() {
		t.Fatalf("erroneously found panic\n")
	}
	reader.ReadStderr([]byte("ic: wrong way"))
	if !reader.FoundPanic() {
		t.Fatalf("failed to find expected 'panic: wrong way'\n")
	}
	if reader.panicType < 0 {
		t.Fatalf("failed to set panic type\n")
	}
	panicType := string(panicHeaders[reader.panicType])
	if panicType != "panic:" {
		t.Fatalf("incorrect panic type: '%s'\n", panicType)
	}
}

func TestSaveExitCode(t *testing.T) {
	reader := NewOutputReader()
	// setup since no process is being run
	reader.group.Add(1)

	reader.AtExit(447)
	if reader.exitCode != 447 {
		t.Fatalf("incorrect process exit code: %d", reader.exitCode)
	}
}
