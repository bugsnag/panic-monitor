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
            | case                  | message                                        | fixture            |
            | explicit panic        | panic: PANIQ!                                  | panic.json         |
            | concurrent read/write | fatal error: concurrent map read and map write | map-readwrite.json |
