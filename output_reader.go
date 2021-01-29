package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	l "github.com/bugsnag/proc-launcher/launcher"
)

const minBufferLen int = 16

type outputReader struct {
	buffer    *bytes.Buffer
	panicType int
	group     *sync.WaitGroup
	exitCode  int
}

func NewOutputReader() *outputReader {
	reader := &outputReader{
		new(bytes.Buffer),
		-1,
		&sync.WaitGroup{},
		0,
	}
	return reader
}

func (reader *outputReader) FoundPanic() bool {
	return reader.panicType != -1
}

func (reader *outputReader) ReadStderr(contents []byte) {
	reader.buffer.Write(contents)
	if !reader.FoundPanic() { // we have yet to find a panic
		for index, header := range panicHeaders {
			location := bytes.Index(reader.buffer.Bytes(), header)
			if location != -1 {
				reader.panicType = index
				reader.buffer.Next(location)
				break
			}
		}
		if !reader.FoundPanic() && reader.buffer.Len() > minBufferLen {
			// truncate all but a buffer long enough to contain an incomplete
			// panic header
			reader.buffer.Truncate(minBufferLen)
		}
	}
}

func (reader *outputReader) AtExit(code int) {
	reader.exitCode = code
	reader.group.Done()
}

func (reader *outputReader) runProcess(args ...string) error {
	// Disable panic handler in the child process
	os.Setenv("BUGSNAG_DISABLE_PANIC_HANDLER", "1")
	// Expose panic handler to child process
	os.Setenv("BUGSNAG_PANIC_MONITOR", "1")
	reader.group.Add(1)
	launcher := l.New(args...)
	launcher.InstallPlugin(reader)
	if err := launcher.Start(); err != nil {
		return fmt.Errorf("Failed to launch process: %v\n", err)
	}
	if err := launcher.Wait(); err != nil {
		return fmt.Errorf("Failed to await process: %v\n", err)
	}
	reader.group.Wait()
	return nil
}

func (reader *outputReader) detectedPanic() (*uncaughtPanic, error) {
	if !reader.FoundPanic() {
		return nil, fmt.Errorf("No panic detected")
	}
	return parsePanic(string(reader.buffer.Bytes()))
}
