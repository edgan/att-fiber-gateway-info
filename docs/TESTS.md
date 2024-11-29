# Tests
These tests require access to a gateway with the exception of the reset tests.
The reset tests just print the question to allow checking the output.

```
# Run all tests
scripts/tests.sh

# Run all tests
scripts/tests.sh all

# Run tests for actions that require login, but we want to test failure
scripts/tests.sh login_failure

# Run tests for actions that do not require a login
scripts/tests.sh nologin

# Run tests for actions that are known to have metrics, and the -allmetrics flag
scripts/tests.sh metrics

# Run tests for actions that require login, but not the reset or restart actions
scripts/tests.sh login

# Run tests for the reset or restart actions
scripts/tests.sh reset

# Output all commands used by tests.sh for debugging
DEBUG=true scripts/tests.sh
```
