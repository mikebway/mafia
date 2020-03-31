// Package cmd utilizes the github.com/spf13/cobra module to parse the
// mafia command line arguements and display usage information.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

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
	Long: `Given an token-code obtained from an MFA device, establishes temporary AWS 
credentials to match a user identity defined in the $HOME/.aws/crdentials file.`,

	// SilenceUsage:  true, // Only display help when explicitly requested, not on error
	SilenceErrors: true, // Only display errors once (helpful when using RunE rathr than Run)

	// RunE is called after the command line has been successfully parsed if no sub-command
	// has been specified. The 'E' indicates that an error (or nil) shall be returned; this
	// cariation of Run is chosen to facilitiate unit testing.
	RunE: func(cmd *cobra.Command, args []string) error {

		// There must be an MFA code
		if len(args) != 1 {
			return errors.New("An MFA token must be provided")
		}

		fmt.Printf("Root command was executed with MFA token: %s\n", args[0])
		return nil
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
