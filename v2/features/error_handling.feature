Feature: Printing information about what went wrong

    Scenario: Running without an API key set
        When I crash the app using no-op
        Then the monitor process exited with an error
        Then the following messages were printed to stderr:
            | Failed to launch monitor: Missing required $BUGSNAG_API_KEY |
            | panic-monitor                                               |
            | EXECUTABLE [EXECUTABLE args]                                |
        And 0 requests were received

    Scenario: Running with an invalid API key set
        When I set the API key to "some-invalid-key"
        When I crash the app using no-op
        Then the monitor process exited with an error
        Then the following messages were printed to stderr:
            | $BUGSNAG_API_KEY must be a 32-character hexadecimal value |
            | panic-monitor                                             |
            | EXECUTABLE [EXECUTABLE args]                              |
        And 0 requests were received

    Scenario: Running with an unknown program
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"
        And I run the monitor with arguments "./unknown-program-name"
        Then the monitor process exited with an error
        And "Failed to run program: " was printed to stderr
        And 0 requests were received

    Scenario: Running without specifying a program
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"
        And I run the monitor with arguments ""
        Then the monitor process exited with an error
        Then the following messages were printed to stderr:
            | No program specified         |
            | panic-monitor                |
            | EXECUTABLE [EXECUTABLE args] |
        And 0 requests were received

    Scenario: Debug logging for failed panic detection
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"
        When I set "DEBUG" to "1" in the environment
        And I run the monitor with:
            | bash | -c | echo "pancake:" >&2 |
        Then the following messages were printed to stderr:
            | pancake:          |
            | No panic detected |
        And 0 requests were received

    Scenario: Debug logging for invalid panic detection
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"
        When I set "DEBUG" to "1" in the environment
        And I run the monitor with:
            | bash | -c | echo "panic: foo" >&2 |
        Then the following messages were printed to stderr:
            | panic: foo             |
            | could not parse panic: |
        And 0 requests were received
