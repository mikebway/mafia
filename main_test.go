// Unit tests for the main Mafia command line entry point.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package main

// Unit tests for the main command line entry point.

import (
	"testing"

	"github.com/mikebway/mafia/cmd"
	"github.com/stretchr/testify/require"
)

// TestMain provides test coverga for the one code line of the main package
// because there is some OCD in my make up!
func TestMain(t *testing.T) {

	// Ensure that we are in a sweet and innocent state and then set up the
	// command line arguments to look like we asked for help and get us
	// the output buffer for us to observe later
	buf := cmd.PrepForExecute()

	// Invoke the main function
	main()

	// Collect the buffered output
	output := buf.String()
	require.Greater(t, len(output), 0, "There should have been some output!!")
	require.Contains(t, output, "mafia token-code [flags]", "The command line help should have been output")
}
