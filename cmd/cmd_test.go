// Unit tests for the cmd package - the Mafia command line parser.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package cmd

// Unit tests for the Cobra command line parsers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/mikebway/mafia/creds"
	"github.com/mikebway/mafia/mfile"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

const (
	// We create (and recreate) a dummy AWS credentials file that we cna control and
	// observe without damaging any geniune article
	fakeCredentialsFilePath = "./credentials.test"

	// The AWS access ID value that we shall populate the fake credentials file with
	fakeAccessKeyID = "FAKE_ACCESS_KEY_ID"

	// The AWS access secret value that we shall populate the fake credentials file with
	fakeSecretAccessKey = "FAKE_SECRET_ACCESS_KEY"

	// The AWS MFA device serial number that we sometimes populate the fake credentials file with
	fakeMFADeviceID = "arn:aws:iam::999999999999:mfa/fake"
)

var (
	// Fake values to be returned by mocked creds.GetSessionCredentials(..)
	accessKey  = "key"
	secret     = "secret"
	token      = "token"
	expiration time.Time

	// Fake sts.GetSessionTokenOutput structure to be returned by mocked creds.GetSessionCredentials(..)
	getSessionTokenOutput = &sts.GetSessionTokenOutput{
		Credentials: &sts.Credentials{
			AccessKeyId:     &accessKey,
			SecretAccessKey: &secret,
			SessionToken:    &token,
			Expiration:      &expiration,
		},
	}
)

// Initialization block
func init() {
	// When running unit test on the command line parser, signal that the actual operations
	// should not be executed, only the parsing.
	unitTesting = true
}

// executeCommand invokes Execute() while capturing its output
// to return for analysis. Any error will have been collected in the
// executeError package global.
func executeCommand(args ...string) string {

	// Clear the ground back to virgin forest and prepare a buffer to capture
	// the execution output
	buf := prepForExecute(args)

	// Execute the command, collecting outout in our buffer
	Execute()

	// Return the output as a string
	return buf.String()
}

// TestBareCommand examines the case where no parameters are provided
func TestBareCommand(t *testing.T) {

	// Run a blank command
	output := executeCommand()

	// We should a complete usage / help dump
	require.Nil(t, executeError, "there should not have been an error: ", executeError)
	require.Contains(t, output,
		"mafia token-code [flags]",
		"Expected usage display")
}

// TestHelpCommand examines the case where no parameters are provided
func TestHelpCommand(t *testing.T) {

	// Run a blank command
	output := executeCommand("help")

	// We should a complete usage / help dump
	require.Nil(t, executeError, "there should not have been an error: ", executeError)
	require.Contains(t, output,
		"mafia token-code [flags]",
		"Expected usage display")
}

// TestBadToken examines the case where an invalid MFA code is provided but nothing else
func TestBadToken(t *testing.T) {

	// Run the command with a bad token value
	output := executeCommand("123")

	// We should have a subcommand required command and a complete usage dump
	require.NotNil(t, executeError, "there should have been an error")
	require.Condition(t, func() bool {
		return checkForExpectedSTSCallFailure(executeError)
	}, "Error should have complained about nonexistent credentials file or invalid MFA token length")

	require.Empty(t, output, "Output for an error condition should have been empty")
}

// TestDisplayHappyPath uses mocking of lower level Mafia packages to prove that the
// command orchestration will successfully display credentials obtained from AWS
// (except we won't actually have called AWS).
func TestDisplayHappyPath(t *testing.T) {

	// Wash the faces of our muddy children before we leave the function
	defer resetChildPackages()

	// Configure our child packages to pretend and return happy answers rather than
	// call AWS or read the real AWS configuration file
	mockChildPackages()

	// Run the command with a random token value that does not matter because we won't actually
	// be calling AWS and so it won't be able to object
	output, stdout := executeCommandCapturingStdout("123456")

	// There should have been no error and no help output
	require.Nil(t, executeError, "there should not have been an error: ", executeError)
	require.Empty(t, output, "there should not have been any help output: %s", output)

	// The stdout capure should contain the environment variables form and the ready-to-paste
	// into credentials file form.
	require.Contains(t, stdout, "export AWS_ACCESS_KEY_ID=key")
	require.Contains(t, stdout, "export AWS_SECRET_ACCESS_KEY=secret")
	require.Contains(t, stdout, "export AWS_SESSION_TOKEN=token")
	require.Contains(t, stdout, "aws_access_key_id = key")
	require.Contains(t, stdout, "aws_secret_access_key = secret")
	require.Contains(t, stdout, "aws_session_token = token")
}

// TestSaveHappyPath uses mocking of lower level Mafia packages to prove that the
// command orchestration will successfully save seession credentials obtained from
// AWS (except we won't actually have called AWS) to a fake AWS credentials file.
func TestSaveHappyPath(t *testing.T) {

	// Wash the faces of our muddy children before we leave the function
	defer resetChildPackages()

	// Configure our child packages to pretend and return happy answers rather than
	// call AWS or read the real AWS configuration file
	mockChildPackages()

	// Run the command with a random token value that does not matter because we won't actually
	// be calling AWS and so it won't be able to object
	output, stdout := executeCommandCapturingStdout("123456", "--save")

	// There should have been no error and no help output
	require.Nil(t, executeError, "there should not have been an error: ", executeError)
	require.Empty(t, output, "there should not have been any help output: %s", output)

	// The stdout capure should contain the environment variables form and the ready-to-paste
	// into credentials file form.
	require.Contains(t, stdout, "Session credentials saved to file")
}

// TestPrepForExecute bumps code coverage by looking at a test prep function that
// would only be otherwise called from the main package test ... which would not
// show in the coverage numbers for this package.
func TestPrepForExecute(t *testing.T) {

	// Start by messing things up so that we can tell that PrepForExecute(..) did something
	unitTesting = false
	executeError = errors.New("this error should get cleared")

	// Call our target function with some parameters that nothing else uses
	buf := PrepForExecute("TestPrepForExecute")

	// Execute the command, collecting output in our buffer
	Execute()

	// Convert the buffer to a string
	output := buf.String()

	// Check everything out
	require.NotNil(t, executeError, "there should have been an error")
	require.Condition(t, func() bool {
		return checkForExpectedSTSCallFailure(executeError)
	}, "Error should have complained about nonexistent credentials file or invalid MFA token length")
	require.Empty(t, output, "Output for an error condition should have been empty")
}

// executeCommandCapturingStdout intercepts stdout and runs the Mafia command with the given
// MFA token. Two strings are returned, the first is the direct output of executeCommand(..)
// and the second is the the captured stdout text.
//
// If somethings goes unfixably wrong with stdout capture, the test run will be aborted altogether.
func executeCommandCapturingStdout(args ...string) (string, string) {

	// We substitute our own pipe for stdout to collect the terminal output
	// but must be careful to always restore stadt and close the pripe files.
	originalStdout := os.Stdout
	readFile, writeFile, err := os.Pipe()
	if err != nil {
		fmt.Printf("Could not capture stdout: %s", err.Error())
		os.Exit(1)
	}

	// Be careful to both put stdout back in its proper place, and restore any
	// tricks that we played on our child packages to get them to cooperate in our testing.
	defer func() {

		// Restore stdout piping
		os.Stdout = originalStdout
		writeFile.Close()
		readFile.Close()
	}()

	// Set our own pipe as stdout
	os.Stdout = writeFile

	// Run the command with a random token value that does not matter because we won't actually
	// be calling AWS and so it won't be able to object
	output := executeCommand(args...)

	// Restore stdout and close the write end of the pipe so that we can collect the output
	os.Stdout = originalStdout
	writeFile.Close()

	// Gather the output into a byte buffer
	outputBytes, err := ioutil.ReadAll(readFile)
	if err != nil {
		fmt.Printf("Failed to read pipe for stdout: : %s", err.Error())
		os.Exit(1)
	}

	// Return the executeCommand output and stdout
	return output, string(outputBytes)
}

// mockChildPackages tricks the kids into behaving the way that we want them to,
// reading the AWS credentials file that we feed them and calling a fake wrapper
// to the AWS STS GetSessionToken(..) that we control.
func mockChildPackages() {

	// Fake an AWS credentials file so that the mfile package will nehave as if it is happy
	setFakeCredentials()

	// Fake out the creds package into using an apparently credentials response from AWS
	creds.SetGetSessionTokenFunc(func(awsService *sts.STS, input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
		return getSessionTokenOutput, nil
	})

}

// setFakeCredentials populates a fake AWS credentials file in the current
// working directory, with an MFA device serial number / ID. The mfile package
// globals are then manipulated such that this fake file will be used any
// future test execution.
func setFakeCredentials() {

	// Start with an empty configuration file content structure
	cfg := ini.Empty()

	// Populate it with the basics
	defaultSection, err := cfg.NewSection(mfile.DefaultSectionName)
	defaultSection.NewKey(mfile.AccessKeyIDKey, fakeAccessKeyID)
	defaultSection.NewKey(mfile.SecretAccessKeyKey, fakeSecretAccessKey)
	defaultSection.NewKey(mfile.MfaDeviceIDKey, fakeMFADeviceID)

	// Write the file
	err = cfg.SaveTo(fakeCredentialsFilePath)

	// That really should not faile to wrote, but if it did abort the tests
	// cos nothing will work after this
	if err != nil {
		fmt.Printf("\nFailed to write fake credentials: %s\n", err.Error())
		os.Exit(999)
	}

	// All looks good - trick the package into using the fake file we just wrote
	mfile.OverrideDefaultCredentialsFilepath(fakeCredentialsFilePath)
}

// resetChildPackages clears any changes that we made to the normal workings of lower
// level Mafia packages to make these tests possible.
func resetChildPackages() {

	// Wash the faces of both dirty kids
	creds.ResetPackageDefaults()
	mfile.ResetPackageDefaults()
}

// checkForExpectedSTSCallFailure checks to see whether one of the expected error conditions occurred
// when calling the AWS STS GetSessionToken(..) method.
func checkForExpectedSTSCallFailure(err error) bool {

	// We expect and error of one or two types ...
	if err != nil ||
		strings.Contains(err.Error(), "Could not read from credentials file") ||
		strings.Contains(err.Error(), "minimum field size of 6") {

		// We saw an error and it was either that there was no AWS credentials file (as
		// will happen with a GitHub action running a build and tests) or that the
		// token parameter was too short to be valid. All is good.
		return true
	}

	// Bummer - this is not what we expected?!?
	return false
}
