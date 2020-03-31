package cmd

// Unit tests for the Cobra command line parsers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
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

	// We should have a subcommand required command and a complete usage dump
	require.NotNil(t, executeError, "there should have been an error")
	require.Equal(t, "An MFA token must be provided", executeError.Error(), "Expected MFA token required error")
	require.Contains(t, output,
		"mafia token-code [flags]",
		"Expected usage display")
}

// TestSimpleOKCommand examines the case where and MFA code is provided but nothing else
func TestSimpleOKCommand(t *testing.T) {

	// Run the command with a token value
	output := executeCommand("123456789")

	// We should have a subcommand required command and a complete usage dump
	require.Nil(t, executeError, "there should have been no error: %v", executeError)
	require.Contains(t, "Root command was executed with MFA token: 123456789", output, "Output not as expected")
}

// TestPrepForExecute bumps code coverage by looking at a test prep function that
// would only be called otherwise from the main package test ... which would not
// show in the coverage numbers for this package.
func TestPrepForExecute(t *testing.T) {

	// Start by messing things up so that we can tell that PrepForExecute(..) did something
	unitTesting = false
	executeError = errors.New("this error should get cleared")

	// Call our target function with some parameters that nothing else uses
	buf := PrepForExecute("TestPrepForExecute")

	// Execute the command, collecting outout in our buffer
	Execute()

	// Convert the buffer to a string
	output := buf.String()

	// Check everything out
	require.Nil(t, executeError, "there should have been no error: %v", executeError)
	require.Contains(t, "Root command was executed with MFA token: TestPrepForExecute", output, "Output not as expected")
}
