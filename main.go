package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	bugsnag "github.com/bugsnag/bugsnag-go/v2"
)

const (
	Version = "1.0.0"
	Usage = `%[1]s: A Go program launcher which automatically reports panics

Usage:

 %[1]s EXECUTABLE [EXECUTABLE args]
 %[1]s -version

 You must specify the environment variable BUGSNAG_API_KEY with your project
 API key prior to launching the program. For more information, see Environment
 Variables.

 Example: %[1]s ./my-app --my-app-flag

Flags:

-version   The version of %[1]s

Environment Variables:

 Override the default values in the monitor through environment variables.

 BUGSNAG_API_KEY               (required) Your API key, available on the
                               Bugsnag dashboard
 BUGSNAG_APP_TYPE              Application component, like a router or queue
 BUGSNAG_APP_VERSION           Current version of the application
 BUGSNAG_ENDPOINT              Event Server address for Bugsnag On-premise
 BUGSNAG_HOSTNAME              Device hostname
 BUGSNAG_METADATA_*            Additional values which will be reported in the
                               event of a panic. See 'Metadata'.
 BUGSNAG_NOTIFY_RELEASE_STAGES Comma-delimited list of stages to notify in
 BUGSNAG_PROJECT_PACKAGES      Comma-delimited list of Go packages to be
                               considered a part of the application
 BUGSNAG_RELEASE_STAGE         The deployment stage of the application, like
                               "production", "beta", or "staging"
 BUGSNAG_SOURCE_ROOT           The directory where source packages are built
                               and the assumed prefix of package directories

Metadata:

 Add metadata through environment variables prefixed with BUGSNAG_METADATA_.

 The environment variable name after the prefix is expected to be the tab and
 key name, delimited by an underscore.

 Underscores in the the tab and/or key values are replaced with spaces.

 Examples:

  Given these environment variables:

  * BUGSNAG_METADATA_device_KubePod="carrot-delivery-service-beta1 reg3"
  * BUGSNAG_METADATA_device_deployment_area=region5_1

  The following fields would be added to a panic report:

  * KubePod: "carrot-delivery-service-beta1 reg3"
  * deployment area: "region5_1"
`
)

var (
	version = flag.Bool("version", false, "The version of panic-monitor")
	APIKeyMatcher *regexp.Regexp = regexp.MustCompile("^[0-9a-fA-f]{32}$")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), Usage, os.Args[0])
	}
	flag.Parse()
	if *version {
		fmt.Printf("%s\n", Version)
		os.Exit(0)
	}
	if err := configureBugsnag(); err != nil {
		printErr("Failed to launch monitor: %v", err)
		flag.Usage()
		os.Exit(1)
	}
	if len(flag.Args()) < 1 {
		printErr("No program specified")
		flag.Usage()
		os.Exit(1)
	}
	reader := NewOutputReader()
	if err := reader.runProcess(flag.Args()...); err != nil {
		printErr("Failed to run program: %v", err)
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
	os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", fmt.Sprintf(format, args...)))
}

func configureBugsnag() error {
	apiKey, err := validateAPIKey(os.Getenv("BUGSNAG_API_KEY"))
	if err != nil {
		return err
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

	bugsnag.Configure(config)
	return nil
}

func validateAPIKey(key string) (string, error) {
	apiKey := os.Getenv("BUGSNAG_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Missing required $BUGSNAG_API_KEY environment variable\n")
	} else if !APIKeyMatcher.MatchString(apiKey) {
		return "", fmt.Errorf("$BUGSNAG_API_KEY must be a 32-character hexadecimal value")
	}

	return apiKey, nil
}
