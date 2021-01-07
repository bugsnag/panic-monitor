package main

import (
	"bytes"
	"fmt"
	"sync"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
	l "github.com/bugsnag/proc-launcher/launcher"
)

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
				break
			}
		}
	}
}

func (reader *outputReader) AtExit(code int) {
	reader.exitCode = code
	reader.group.Done()
}

func (reader *outputReader) runProcess(args ...string) error {
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

func (reader *outputReader) checkAndSendPanicEvent() {
	if !reader.FoundPanic() {
		return
	}
	contents := string(reader.buffer.Bytes())
	value, err := parsePanic(contents)
	if value != nil {
		bugsnag.Notify(value, bugsnag.HandledState{
			SeverityReason:   bugsnag.SeverityReasonUnhandledPanic,
			OriginalSeverity: bugsnag.SeverityError,
			Unhandled:        true,
		}, bugsnag.ErrorClass{Name: value.typeName})
	} else {
		fmt.Printf("Could not deliver panic due to error: %v\n", err)
	}
}
