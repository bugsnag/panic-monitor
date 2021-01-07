# panic-monitor

An executable which launches a program, awaiting its exit. In the event of
a go panic, it is automatically reported to Bugsnag.

## Usage

Set your Bugsnag API key as the environment variable `BUGSNAG_API_KEY`, then
launch your program

```sh
panic-monitor PROGRAM [PROGRAM arguments]
```

### Configuration

Override the default values in the monitor through environment variables:

| Key                        | Value |
|----------------------------|-------|
| BUGSNAG\_API\_KEY          | **(required)** Your API key, available on the Bugsnag dashboard |
| BUGSNAG\_APP\_TYPE         | Application component, like a router, mailer, or queue|
| BUGSNAG\_APP\_VERSION      | Current version of the application |
| BUGSNAG\_ENDPOINT          | Event Server address for Bugsnag On-premise |
| BUGSNAG\_HOSTNAME          | Device hostname |
| BUGSNAG\_PROJECT\_PACKAGES | Comma-delimited list of Go packages to be considered a part of the application |
| BUGSNAG\_RELEASE\_STAGE    | The deployment stage of the application, like "production" or "beta" or "staging" |
| BUGSNAG\_SOURCE\_ROOT      | The directory where source packages are built and the assumed prefix of package directories |

## Examples

Build one of the example crashing apps using `go build`:

```sh
TESTCASE="explicit panic" go build features/fixtures/app
```

Then run it using the monitor:

```sh
BUGSNAG_API_KEY="your-api-key-here" ./panic-monitor ./app
```

## Testing

Run the unit tests via `go test`

The integration tests depend on cucumber, available through Ruby bundler:

```sh
bundle install
bundle exec cucumber
```
