package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

func main() {
	if err := configureBugsnag(); err != nil {
		printErr("Failed to launch monitor: %v\n", err)
		printUsage()
		os.Exit(1)
	}
	if len(os.Args) < 2 {
		printErr("No program specified\n")
		printUsage()
		os.Exit(1)
	}
	reader := NewOutputReader()
	if err := reader.runProcess(os.Args[1:]...); err != nil {
		printErr("Failed to run program: %v\n", err)
		os.Exit(1)
	}

	event, detectionErr := reader.detectedPanic()
	if detectionErr != nil && debugModeEnabled() {
		printErr("%v", detectionErr)
	}
	if event != nil {
		bugsnag.Notify(event, bugsnag.HandledState{
			SeverityReason:   bugsnag.SeverityReasonUnhandledPanic,
			OriginalSeverity: bugsnag.SeverityError,
			Unhandled:        true,
		}, bugsnag.ErrorClass{Name: event.typeName})
	}

	os.Exit(reader.exitCode)
}

func debugModeEnabled() bool {
	return os.Getenv("DEBUG") == "1"
}

func printErr(format string, args ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(format, args...))
}

func printUsage() {
	fmt.Printf("Usage: %s EXECUTABLE [EXECUTABLE args]\n", os.Args[0])
}

func configureBugsnag() error {
	matcher := regexp.MustCompile("^[0-9a-fA-f]{32}$")
	apiKey := os.Getenv("BUGSNAG_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("Missing required $BUGSNAG_API_KEY environment variable\n")
	} else if !matcher.MatchString(apiKey) {
		return fmt.Errorf("$BUGSNAG_API_KEY must be a 32-character hexadecimal value")
	}

	config := bugsnag.Configuration{
		APIKey:              apiKey,
		PanicHandler:        func() {},
		Synchronous:         true,
		AutoCaptureSessions: false,
	}

	if !debugModeEnabled() {
		// mute warnings from underlying bugsnag lib unless in debug mode
		config.Logger = &logger{}
	}

	if stage := os.Getenv("BUGSNAG_RELEASE_STAGE"); stage != "" {
		config.ReleaseStage = stage
	}
	if app_version := os.Getenv("BUGSNAG_APP_VERSION"); app_version != "" {
		config.AppVersion = app_version
	}
	if hostname := os.Getenv("BUGSNAG_HOSTNAME"); hostname != "" {
		config.Hostname = hostname
	}
	if sourceRoot := os.Getenv("BUGSNAG_SOURCE_ROOT"); sourceRoot != "" {
		config.SourceRoot = sourceRoot
	}
	if appType := os.Getenv("BUGSNAG_APP_TYPE"); appType != "" {
		config.AppType = appType
	}
	if stages := os.Getenv("BUGSNAG_NOTIFY_RELEASE_STAGES"); stages != "" {
		config.NotifyReleaseStages = strings.Split(stages, ",")
	}
	if packages := os.Getenv("BUGSNAG_PROJECT_PACKAGES"); packages != "" {
		config.ProjectPackages = strings.Split(packages, ",")
	}
	if endpoint := os.Getenv("BUGSNAG_ENDPOINT"); endpoint != "" {
		config.Endpoints = bugsnag.Endpoints{
			Notify:   endpoint,
			Sessions: "",
		}
	}

	bugsnag.OnBeforeNotify(func(event *bugsnag.Event, config *bugsnag.Configuration) error {
		for _, value := range os.Environ() {
			key, value, err := parseEnvironmentPair(value)
			if err != nil {
				continue
			}
			if keypath, err := parseMetadataKeypath(key); err == nil {
				addMetadata(event, keypath, value)
			}
		}

		return nil
	})

	bugsnag.Configure(config)
	return nil
}
