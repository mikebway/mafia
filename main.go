// A command line utility supporting MFA authentication to obtain
// temporary AWS session credentials wher second factor auhentication
// is required.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package main

import "github.com/mikebway/mafia/cmd"

// Command line entry point.
//
// Cobra based command line parsing does all the work.
func main() {
	cmd.Execute()
}
