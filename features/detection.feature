Feature: Detecting panic-monitor from the child process

    Background:
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"

    Scenario: Changing behavior based on monitor env variable
        When I run the monitor with:
            | bash | -c | echo "monitor: $BUGSNAG_PANIC_MONITOR" |
        Then "monitor: 1" was printed to stdout
