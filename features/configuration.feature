Feature: Configuring Bugsnag

    Background:
        When I set the API key to "035d2472bd130ac0ab0f52715bbdc65d"

    Scenario Outline: Adding content to events through configuration
        When I set "<variable>" to "<value>" in the environment
        And I crash the app using explicit panic
        Then the monitor process exited with an error
        And 1 request was received
        And payload field "<field>" equals "<value>"

        Examples:
            | variable              | value           | field             |
            | BUGSNAG_APP_VERSION   | 1.4.34          | app.version       |
            | BUGSNAG_APP_TYPE      | mailer-daemon   | app.type          |
            | BUGSNAG_RELEASE_STAGE | beta1           | app.releaseStage  |
            | BUGSNAG_HOSTNAME      | dream-machine-2 | device.hostname   |

            | BUGSNAG_METADATA_device.instance      | kube2-33-A | metaData.device.instance      |
            | BUGSNAG_METADATA_framework.version    | v3.1.0     | metaData.framework.version    |
            | BUGSNAG_METADATA_device.runtime_level | 1C         | metaData.device.runtime level |
            | BUGSNAG_METADATA_Carrot               | orange     | metaData.custom.Carrot        |

    Scenario: Configuring project packages
        When I set "BUGSNAG_PROJECT_PACKAGES" to "main,github.com/bugsnag/panic-monitor/features/fixtures/app" in the environment
        And I crash the app using explicit panic
        Then the monitor process exited with an error
        And 1 request was received
        And the payload contains the following in-project stack frames:
            | file     | method        | lineNumber |
            | cases.go | explicitPanic | 17         |
            | main.go  | main          | 11         |

    Scenario: Configuring source root
        When I set "BUGSNAG_SOURCE_ROOT" to the sample app directory
        And I crash the app using explicit panic
        Then the monitor process exited with an error
        And 1 request was received
        And the payload contains the following in-project stack frames:
            | file     | method        | lineNumber |
            | cases.go | explicitPanic | 17         |
            | main.go  | main          | 11         |

    Scenario: Delivering events filtering through notify release stages
        When I set "BUGSNAG_NOTIFY_RELEASE_STAGES" to "prod,beta" in the environment
        When I set "BUGSNAG_RELEASE_STAGE" to "beta" in the environment
        And I crash the app using explicit panic
        Then the monitor process exited with an error

    Scenario: Suppressing events through notify release stages
        When I set "BUGSNAG_NOTIFY_RELEASE_STAGES" to "prod,beta" in the environment
        When I set "BUGSNAG_RELEASE_STAGE" to "dev" in the environment
        And I crash the app using explicit panic
        Then the monitor process exited with an error
        And 0 requests were received
