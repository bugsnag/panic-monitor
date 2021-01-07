# panic-monitor

An executable which launches a program, awaiting its exit. In the event of
a go panic, it is automatically reported to Bugsnag.

## Usage

Set your Bugsnag API key as the environment variable `BUGSNAG_API_KEY`, then
launch your program

```sh
panic-monitor PROGRAM [PROGRAM arguments]
```

## Examples

Build one of the example crashing apps using `go build`:

```sh
TESTCASE="explicit panic" go build features/fixtures/app
```

Then run it using the monitor:

```sh
BUGSNAG_API_KEY="your-api-key-here" ./panic-monitor ./app
```
