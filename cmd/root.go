// Package cmd utilizes the github.com/spf13/cobra module to parse the
// mafia command line arguements and display usage information.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/mikebway/mafia/creds"
	"github.com/mikebway/mafia/mfile"
	"github.com/spf13/cobra"
)

var (
	unitTesting  = false // Set to true when running unit tests
	executeError error   // The error value obtained by Execute(), captured for unit test purposes

	saveCredentials = false // True if update the $HOME/.aws/credentials file with the session credentionals obtained
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mafia token-code",
	Short: "Establishes temporary AWS credentials where MFA authentication codes are required",
	Long: `
Given a token/number obtained from an MFA device, establishes temporary AWS 
credentials to match a user identity defined in the ~/.aws/crdentials file.

Before running, you must add your MFA serial number to the [default] section of
the ~/.aws/crdentials file, alongside the aws_access_key_id and 
aws_secret_access_key values, as follows:

   mfa_device_id = arn:aws:iam::999999999999:mfa/jane

Replacing 999999999999 with your account number, and jane with your username.
`,

	SilenceUsage:  true, // Only display help when explicitly requested, not on error
	SilenceErrors: true, // Only display errors once (helpful when using RunE rathr than Run)

	// RunE is called after the command line has been successfully parsed if no sub-command
	// has been specified. The 'E' indicates that an error (or nil) shall be returned; this
	// cariation of Run is chosen to facilitiate unit testing.
	RunE: func(cmd *cobra.Command, args []string) error {

		// If no MFA code was provided or help was requested, display the help
		if len(args) != 1 || args[0] == "help" {
			return cmd.Help()
		}

		// Do the work!
		credentials, err := fetchSessionCredentials(args[0])
		if err != nil {
			return err
		}

		// If we are to save the credentials ...
		if saveCredentials {

			// Try to the save the credentials
			err = saveSessionCredentials(credentials)

			// If that worked, give the user a comfort signal
			fmt.Println("Session credentials saved to file")

		} else {

			// Not saving the credentials so show them in stdout
			displaySessionCredentials(credentials)
		}

		// All done - maybe not successfully; either way return the rror value that we have
		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if executeError = rootCmd.Execute(); executeError != nil {
		fmt.Println(executeError)
		if !unitTesting {
			os.Exit(1)
		}
	}
}

// Load time initialization - called automatically
func init() {

	// Initialize the flags that apply to the root command and, potentially, to subcommands
	initRootFlags()
}

// initRootFlags is called from init() to define the flags that apply to the root
// command, and might be inherited by its subcommands. It is defined separately from
// init() so that it can be invoked by unit tests when they need to reset the playing field.
func initRootFlags() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().BoolVar(&saveCredentials, "save", false, "save the obtained credentials to the .aws/credentials file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//  rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// ============================================================================
// The following ar provided to support unit tests. In particular, they allow
// the tests for the main package to ensure that the environment is reset
// from any prior tests.
// ============================================================================

// PrepForExecute is FOR UNIT TESTING ONLY. It is is intended to be invoked by
// the main package's tests to clear any previous test debris and capture output
// in the returned buffer. It accepts a variable number of string values to
// be set as arguments for command execution.
func PrepForExecute(args ...string) *bytes.Buffer {

	// Have our string slice accepting sibling do all the work
	return prepForExecute(args)
}

// prepForExecute is a package scoped function accepting a string array of
// command line argument values. It resets the session environment and
// establishs a a buffer to capture the output of the subsequent command
// execution.
func prepForExecute(args []string) *bytes.Buffer {

	// Ensure that we are in a sweet and innocent state
	resetCommand()

	// Set the arguments and invoke the normal Execute() package entry point
	rootCmd.SetArgs(args)

	// Arrange to collect the output in a buffer
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	// Return the, as yet, empty buffer
	return buf
}

// resetCommand clears both command specific parameter values and
// global ones so that tests can be run in a known "virgin" state.
func resetCommand() {

	// Tell the world that we are unit testing and not to get
	// too carried away
	unitTesting = true

	// Clear and then re-initialize all the flags definitions
	rootCmd.ResetFlags()
	initRootFlags()
}

// fetchSessionCredentials orchestrates the work of obtaining, displaying, and
// potentially saving AWS session credentials to the  ~/.aws/credentials file.
func fetchSessionCredentials(mfaToken string) (*creds.SessionCredentials, error) {

	// Obtain the MFA device ID / serial number as defined by AWS
	mfaDeviceID, err := mfile.GetMFADeviceID()
	if err != nil {
		return nil, err
	}

	// Ask AWS for the credentials and return what we get
	return creds.GetSessionCredentials(mfaDeviceID, mfaToken, 3600)
}

// displaySessionCredentials shows the, you guessed it, session credentials on stdout.
// The display is given twice, once formated for use as environment variables and
// once ready to copy-nd-paste into the  ~/.aws/credentials file.
func displaySessionCredentials(credentials *creds.SessionCredentials) {

	// Display the results in a form that can be copy-and-pasted to set as environment variables
	fmt.Printf("\nEnvironment Variables\n\n")
	fmt.Printf("export AWS_ACCESS_KEY_ID=%s\n", *credentials.AccessKeyID)
	fmt.Printf("export AWS_SECRET_ACCESS_KEY=%s\n", *credentials.SecretAccessKey)
	fmt.Printf("export AWS_SESSION_TOKEN=%s\n", *credentials.SessionToken)
	fmt.Println("history -c # clear shell history immediatly after setting secrets")

	// Display the results in a form that can be copy-and-pasted to set as environment variables
	fmt.Printf("\nTo paste into ~/.aws/credentials\n\n")
	fmt.Println("[default-session]")
	fmt.Printf("aws_access_key_id = %s\n", *credentials.AccessKeyID)
	fmt.Printf("aws_secret_access_key = %s\n", *credentials.SecretAccessKey)
	fmt.Printf("aws_session_token = %s\n", *credentials.SessionToken)
	fmt.Println()
}

// saveSessionCredentials attempts to svae the obtained session credentials to the
// ~/.aws/credentials file.
func saveSessionCredentials(credentials *creds.SessionCredentials) error {

	return mfile.SaveSessionCredentials(credentials.AccessKeyID, credentials.SecretAccessKey, credentials.SessionToken)
}
