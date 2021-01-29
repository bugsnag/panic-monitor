Feature: Reporting fatal panics

    Background:
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"

    Scenario Outline: Fatal panics from varying causes
        When I crash the app using <case>
        Then the monitor process exited with an error
        And 1 request was received
        And "<message>" was printed to stderr
        And I receive an error event matching <fixture>

        Examples:
            | case                  | message                                        | fixture             |
            | explicit panic        | panic: PANIQ!                                  | panic.json          |
            | goroutine             | panic: at the disco?                           | goroutine.json      |
            | concurrent read/write | fatal error: concurrent map read and map write | map-readwrite.json  |
            | nil goroutine         | fatal error: go of nil func value              | nil-goroutine.json  |
            | fake panic then real  | panic: REAL PANIC!                             | garbage-panic.json  |
            | stack exhaustion      | fatal error: stack overflow                    | stack-overflow.json |
            | array overflow        | panic: runtime error: index out of range       | array-overflow.json |
            | nil pointer deref     | runtime error: invalid memory address or nil pointer dereference | nil-pointer.json    |
            | bad reflect swap      | panic: reflect: call of Swapper on string Value                  | reflect-swap.json   |

    Scenario: Avoid double-reporting panics
        When I crash the bugsnag-app using explicit panic
        And I wait for 2 seconds
        Then the monitor process exited with an error
        And 1 request was received
        And "oh no!" was printed to stderr
        And I receive an error event matching oh-no-panic.json
