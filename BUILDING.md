# Building
## Compiling
To build for only system's platform:
```
go build
```

To binaries for all supported combinations of operating systems and architectures:
```
scripts/builds.sh
```

## Builds
See the `.go_builds` file for the list of supported combinations of operating
systems and architectures.

## Tests
These tests require access to a gateway with the exception of the reset tests.
The reset tests just print the question to allow checking the output.

```
# Run all tests
scripts/tests.sh

# Run all tests
scripts/tests.sh all

# Run tests for actions that do not require a login
scripts/tests.sh nologin

# Run tests for actions that are known to have metrics, and the -allmetrics flag
scripts/tests.sh metrics

# Run tests for actions that require login, but not the reset or restart actions
scripts/tests.sh login

# Run tests for the reset or restart actions
scripts/tests.sh reset
```
