# Mafia - A CLI Utility Supporting AWS MFA Authentication

**Mafia** may be used as a subsitute for the `aws sts get-session-token` command, offering the
advantages of

* Not having to remember and type out your MFA device serial number.

* Displaying the obatained session credentials in copy-and paste ready
  forms for both:

  * Setting environment variables
  * Adding to your `$HOME/.aws/credentials` file

* Optionally, writing the session credentials to the `$HOME/.aws/credentials` file
  for you.

## Usage

Running the `mafia` utility with any of an argument of `help`, or the `--help` or `-h`
flags will display the following usage information:

```text
Given a token/number obtained from an MFA device, establishes temporary AWS
credentials to match a user identity defined in the ~/.aws/crdentials file.

Before running, you must add your MFA serial number to the [default] section of
the ~/.aws/crdentials file, alongside the aws_access_key_id and
aws_secret_access_key values, as follows:

   mfa_device_id = arn:aws:iam::999999999999:mfa/jane

Replacing 999999999999 with your account number, and jane with your username.

Usage:
  mafia token-code [flags]

Flags:
  -h, --help   help for mafia
      --save   save the obtained credentials to the .aws/credentials file
```

Note especially the need to declare your MFA device ID / serial number in the
`$HOME/.aws/credentials` file.

## What's Missing

* A flag to specifiy something other than the default credentials in the
`$HOME/.aws/credentials` file.

* A flag to specify the name and path of the credentials file, other than the
default `$HOME/.aws/credentials` location.

* Support for Microsoft Windows users, where credential file paths are
specified differently. The author is unlikely to get to that since they don't 
have a Windows system to build and test with.

## Unit / Integration Testing

The unit tests are really more like integration tests in that they will invoke
AWS API calls though successful calls are only achieved through mocking.

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

Unit test coverage should be kept above 90% by line for all packages.
