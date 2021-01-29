Feature: Detecting panic-monitor from the child process

    Background:
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"

    Scenario: Changing behavior based on monitor env variable
        When I run the monitor with:
            | bash | -c | echo "monitor: $BUGSNAG_PANIC_MONITOR" |
        Then "monitor: 1" was printed to stdout

    Scenario: Disabling panic handler in the child process by default
        When I run the monitor with:
            | bash | -c | echo "disabled: $BUGSNAG_DISABLE_PANIC_HANDLER" |
        Then "disabled: 1" was printed to stdout
