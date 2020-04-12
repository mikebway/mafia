package mfile

// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
//
// See doc.go for other overall package documentation. This file contains
// package methods related to reading the AWS credentials file and also
// defines the package constants and global variables.

import (
	"fmt"
	"os"
	"os/user"

	"github.com/go-ini/ini"
)

const (
	// The name of the default section in the AWS configuration file
	defaultSectionName = "default"

	// The key name of the AWS accesss key ID field within a configuration file section
	accessKeyIDKey = "aws_access_key_id"

	// The key name of the AWS secret access key field within a configuration file section
	secretAccessKeyKey = "aws_secret_access_key"

	// The key name of any MFA authenticated temporary session token field within a configuration file section
	sessionTokenKey = "aws_session_token"

	// The key name of the MFA device ID field within a configuration file section
	mfaDeviceIDKey = "mfa_device_id"

	// Suffix appended to the non-session section name to name the correseponding
	// MHF authenticated session credentials section
	sessionSectionSuffix = "-session"
)

var (
	// What the name says, filled in at load time. As a global variable, this can be
	// overriden by unit tests to better control outcomes.
	defaultCredentialsFilePath string
)

// Load time initialization
func init() {

	// Configure the location fo the AWS credentials file and other initial state
	// assumed by this package.
	ResetPackageDefaults()
}

// GetMFADeviceID attempts to find an MFA device ID in the default section of
// the AWS credentials file, returing either the ID or an error.
func GetMFADeviceID() (string, error) {
	return GetMFADeviceIDFromFile(defaultCredentialsFilePath)
}

// GetMFADeviceIDFromFile attempts to find an MFA device ID in the default section of
// the given AWS credentials file, returing either the ID or an error.
func GetMFADeviceIDFromFile(filepath string) (string, error) {

	// Load the file
	cfg, err := ini.Load(filepath)
	if err != nil {
		return "", fmt.Errorf("Could not read from credentials file %s: %v", filepath, err)
	}

	// Fetch the default section - if there is one
	defaultSection := cfg.Section(defaultSectionName)
	if len(defaultSection.Keys()) == 0 {
		return "", fmt.Errorf("%s section not found or empty in %s", defaultSectionName, filepath)
	}

	// Fetch the MFA device ID entry - if there is one
	key := defaultSection.Key(mfaDeviceIDKey)
	if len(key.Value()) == 0 {
		return "", fmt.Errorf("%s key not found in default section of %s", mfaDeviceIDKey, filepath)
	}

	// Return the value of the key
	return key.String(), nil
}

// OverrideDefaultCredentialsFilepath is intended for use by unit tests that need to
// manage the behavior of this package when loading and saving to the 'default'
// AWS credentials file, protecting the real file from being damaged ny the tests.
func OverrideDefaultCredentialsFilepath(filepath string) {
	defaultCredentialsFilePath = filepath
}

// ResetPackageDefaults ensures that the package is in its proper default state, ready
// to go to work. This is used when the package is first loaded but also by unit tests
// needing to restore initial conditions after a potentially destructive test run.
func ResetPackageDefaults() {

	// Set the path for the default AWS credentials file
	defaultCredentialsFilePath = getDefaultCredentialsFilepath()
}

// getDefaultCredentialsFilepath obtains the home directory of the current
// user and forms the full path to the default AWS credentials file from that.
func getDefaultCredentialsFilepath() string {

	// Ask the OS for information about the current user. This should never fail
	// but if it does, barf and kill the program here and now.
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to obtain credentials information: %v\n", err)
		os.Exit(1)
	}

	// Configure the default AWS credentials path
	return usr.HomeDir + "/.aws/credentials"
}
