# panic-monitor

An executable which launches a program, awaiting its exit. In the event of
a go panic, it is automatically reported to Bugsnag.

## Installation

```sh
go install github.com/bugsnag/panic-monitor@latest
```

## Usage

Set your Bugsnag API key as the environment variable `BUGSNAG_API_KEY`, then
launch your program

```sh
panic-monitor PROGRAM [PROGRAM arguments]
```

### Detection

Detect when a program was launched via panic-monitor by the presence of the
environment variable `BUGSNAG_PANIC_MONITOR`.

### Configuration

Override the default values in the monitor through environment variables:

| Key                              | Value |
|----------------------------------|-------|
| BUGSNAG\_API\_KEY                | **(required)** Your API key, available on the Bugsnag dashboard |
| BUGSNAG\_APP\_TYPE               | Application component, like a router, mailer, or queue|
| BUGSNAG\_APP\_VERSION            | Current version of the application |
| BUGSNAG\_NOTIFY\_ENDPOINT        | Event Server address for Bugsnag On-premise |
| BUGSNAG\_HOSTNAME                | Device hostname |
| BUGSNAG\_NOTIFY\_RELEASE\_STAGES | Comma-delimited list of release stages to notify in |
| BUGSNAG\_PROJECT\_PACKAGES       | Comma-delimited list of Go packages to be considered a part of the application |
| BUGSNAG\_RELEASE\_STAGE          | The deployment stage of the application, like "production" or "beta" or "staging" |
| BUGSNAG\_SOURCE\_ROOT            | The directory where source packages are built and the assumed prefix of package directories |

### Custom metadata

Add metadata through environment variables prefixed with `BUGSNAG_METADATA_`.

The environment variable name after the prefix is expected to be the tab and key name,
delimited by an underscore.

Examples:

```sh
BUGSNAG_METADATA_device_KubePod="carrot-delivery-service-beta1 reg3"
BUGSNAG_METADATA_device_deployment_area=region5_1
```

Would add the following metadata to the `device` tab in the event of a panic:

* `KubePod`: `carrot-delivery-service-beta1 reg3`
* `deployment_area`: `region5_1`

## Examples

Build one of the example crashing apps using `go build`:

```sh
go build features/fixtures/app
```

Then run it using the monitor:

```sh
export BUGSNAG_API_KEY="your-api-key-here"
TESTCASE="explicit panic" panic-monitor ./app
```

## Testing

Run the unit tests via `go test`

The integration tests depend on cucumber, available through Ruby bundler:

```sh
bundle install
bundle exec cucumber
```
