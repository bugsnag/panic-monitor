package main

import (
	"fmt"
	"os"
	"regexp"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

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
	reader := NewOutputReader()
	if err := reader.runProcess(os.Args[1:]...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	event, _ := reader.detectedPanic()
	if event != nil {
		bugsnag.Notify(event, bugsnag.HandledState{
			SeverityReason:   bugsnag.SeverityReasonUnhandledPanic,
			OriginalSeverity: bugsnag.SeverityError,
			Unhandled:        true,
		}, bugsnag.ErrorClass{Name: event.typeName})
	}

	os.Exit(reader.exitCode)
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
