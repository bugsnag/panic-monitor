package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sync"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
	"github.com/bugsnag/bugsnag-go/v2/errors"
	l "github.com/bugsnag/proc-launcher/launcher"
)

type extension struct {
	buffer    *bytes.Buffer
	headers   [][]byte
	panicType int
	group     *sync.WaitGroup
	exitCode  int
}

func New() *extension {
	ex := &extension{
		new(bytes.Buffer),
		panicHeaders(),
		-1,
		&sync.WaitGroup{},
		0,
	}
	ex.group.Add(1)
	return ex
}

func (ex *extension) FoundPanic() bool {
	return ex.panicType != -1
}

func (ex *extension) ReadStderr(contents []byte) {
	ex.buffer.Write(contents)
	if !ex.FoundPanic() { // we have yet to find a panic
		for index, header := range ex.headers {
			location := bytes.Index(ex.buffer.Bytes(), header)
			if location != -1 {
				ex.panicType = index
				break
			}
		}
	}
}

func (ex *extension) AtExit(code int) {
	ex.exitCode = code
	ex.group.Done()
}

func main() {
	if err := configureBugsnag(); err != nil {
		fmt.Printf("Failed to launch monitor: %v\n", err)
		printUsage()
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	ex := New()
	ex.runProcess()
	ex.checkAndSendPanicEvent()
	os.Exit(ex.exitCode)
}

func printUsage() {
	fmt.Printf("Usage: %s EXECUTABLE [EXECUTABLE args]\n", os.Args[0])
}

func configureBugsnag() error {
	matcher := regexp.MustCompile("^[0-9a-fA-f]{32}$")
	apiKey := os.Getenv("BUGSNAG_API_KEY")
	if !matcher.MatchString(apiKey) {
		return fmt.Errorf("Missing required $BUGSNAG_API_KEY environment variable\n")
	}

	config := bugsnag.Configuration{
		APIKey:              apiKey,
		PanicHandler:        func() {},
		Synchronous:         true,
		AutoCaptureSessions: false,
		Logger:              &logger{},
	}

	endpoint := os.Getenv("BUGSNAG_ENDPOINT")
	if endpoint != "" {
		config.Endpoints = bugsnag.Endpoints{
			Notify:   endpoint,
			Sessions: "",
		}
	}

	bugsnag.Configure(config)
	return nil
}

func (ex *extension) runProcess() {
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
	ex.group.Wait()
}

func (ex *extension) checkAndSendPanicEvent() {
	if !ex.FoundPanic() {
		return
	}
	// this logic should actually move into this package
	contents := string(ex.buffer.Bytes())
	value, err := errors.ParsePanic(contents)
	if value != nil {
		bugsnag.Notify(value, bugsnag.HandledState{
			SeverityReason:   bugsnag.SeverityReasonUnhandledPanic,
			OriginalSeverity: bugsnag.SeverityError,
			Unhandled:        true,
		})
	} else {
		fmt.Printf("Could not deliver panic due to error: %v\n", err)
	}
}
