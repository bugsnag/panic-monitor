package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/bugsnag/bugsnag-go/errors"
	l "github.com/bugsnag/proc-launcher/launcher"
)

type extension struct {
	buffer  *bytes.Buffer
	headers [][]byte
	panicType int
}

func New() *extension {
	return &extension{
		new(bytes.Buffer),
		[][]byte{
			[]byte("panic:"),
			[]byte("fatal error:"),
		},
		-1,
	}
}

func (ex *extension) FoundPanic() bool {
	return ex.panicType != -1
}

func (ex *extension) ReadStderr(contents []byte) {
	ex.buffer.Write(contents)
	if !ex.FoundPanic() { // we have yet to find a panic
		for index, header := range ex.headers {
			if bytes.Index(ex.buffer.Bytes(), header) != -1 {
				ex.panicType = index
				break
			}
		}
		if !ex.FoundPanic() {
			ex.buffer.Truncate(16)
		}
	}
}

func main() {
	matcher := regexp.MustCompile("^[0-9a-fA-f]{32}$")
	apiKey := os.Getenv("BUGSNAG_API_KEY")
	if !matcher.MatchString(apiKey) {
		fmt.Printf("Set $BUGSNAG_API_KEY environment variable to launch monitor\n")
		fmt.Printf("Usage: %s EXECUTABLE [EXECUTABLE args]\n", os.Args[0])
		os.Exit(1)
	}
	bugsnag.Configure(bugsnag.Configuration{
		APIKey: apiKey,
		PanicHandler: func() {},
		Synchronous: true,
		AutoCaptureSessions: false,
	})
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s EXECUTABLE [EXECUTABLE args]\n", os.Args[0])
		os.Exit(1)
	}
	ex := New()
	launcher := l.New(os.Args[1:]...)
	launcher.InstallPlugin(ex)
	if err := launcher.Start(); err != nil {
		fmt.Printf("failed to launch process: %v\n", err)
		os.Exit(1)
	}
	if err := launcher.Wait(); err != nil {
		fmt.Printf("failed to await process: %v\n", err)
		os.Exit(1)
	}
	if ex.FoundPanic() {
		// this logic should actually move into this package
		value, _ := errors.ParsePanic(string(ex.buffer.Bytes()))
		if value != nil {
			bugsnag.Notify(value, bugsnag.HandledState{
				SeverityReason: bugsnag.SeverityReasonUnhandledPanic,
				OriginalSeverity: bugsnag.SeverityError,
				Unhandled: true,
			})
		}
	}
}
