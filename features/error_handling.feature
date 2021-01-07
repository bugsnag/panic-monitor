Feature: Printing information about what went wrong

    Background:
        Given I build the executable
        And I build a sample app

    Scenario: Running without an API key set
        When I crash the app using no-op
        Then the monitor process exited with an error
        And "Failed to launch monitor: Missing required $BUGSNAG_API_KEY" was printed to stderr
        And "Usage: ./panic-monitor EXECUTABLE [EXECUTABLE args]" was printed to stdout
        And 0 requests were received

    Scenario: Running with an invalid API key set
        When I set the API key to "some-invalid-key"
        When I crash the app using no-op
        Then the monitor process exited with an error
        And "$BUGSNAG_API_KEY must be a 32-character hexadecimal value" was printed to stderr
        And "Usage: ./panic-monitor EXECUTABLE [EXECUTABLE args]" was printed to stdout
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
        And "No program specified" was printed to stderr
        And "Usage: ./panic-monitor EXECUTABLE [EXECUTABLE args]" was printed to stdout
        And 0 requests were received
