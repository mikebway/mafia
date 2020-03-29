# Mafia - A CLI Utility Supporting AWS MFA Authentication

**Mafia** may be used as a subsitute for the `aws sts get-session-token` command, offering the
advantages of

* Not having to remember and type out your MFA device serial number.

* Displaying the obataind session credentials in copy-and paste ready
  forms for both:

  * Setting environment variables
  * Adding to your `$HOME/.aws/credentials` file

* Optionally, writing the session credentials to the `$HOME/.aws/credentials` file
  for you.

## Usage

Running the `mafia` utility with either the `--help` or `-h` flags will display the
following usage information:

```text
Given an token-code obtained from an MFA device, establishes temporary AWS
credentials to match a user identity defined in the $HOME/.aws/crdentials file.

Usage:
  mafia token-code [flags]

Flags:
  -h, --help   help for mafia
      --save   save the obtained credentials to the .aws/credentials file
```

## Unit / Integration Testing

The unit tests are really more like integration tests in that they will invoke
AWS API calls; mocks are not used.

You can run all of the unit tests from the command line and receive a coverage
report as follows:

```bash
go test -cover ./...
```

To ensure that all tests are run, and that none are assumed unchanged for the
cache of a previous run, you may add the `-count=1` flag to required that all
tests are run at least and exactly once:

```bash
go test -cover -count=1 ./...
```