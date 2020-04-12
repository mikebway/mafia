package mfile

// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
//
// See doc.go for package documentation

import (
	"fmt"
	"os"
	"os/user"

	"github.com/go-ini/ini"
)

const (
	// The name of the default section in the AWS configuration file
	defaultSectionName = "default"

	// The key name of the MFA device ID field within a configuration file section
	mfaDeviceIDKey = "mfa_device_id"

	// Suffix appended to the non-session section name to name the correseponding
	// MHF authenticated session credentials section
	sessionSectionSuffix = "-session"
)

// getDefaultCredentialsFilepath obtains the home directory of the current
// user and forms the full path to the default AWS credentials file from that.
func getDefaultCredentialsFilepath() string {

	// Ask the OS for information about the current user. This should never fail
	// but if it does, barf and kill the program here and now.
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Failed to obtain user information: %v\n", err)
		os.Exit(1)
	}

	// Configure the default AWS credentials path
	return usr.HomeDir + "/.aws/credentials"
}

// GetMFADeviceID attempts to find an MFA device ID in the default section of
// the AWS credentials file, returing either the ID or an error.
func GetMFADeviceID() (string, error) {
	return GetMFADeviceIDFromFile(getDefaultCredentialsFilepath())
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
	if defaultSection == nil {
		return "", fmt.Errorf("%s section not found in %s", defaultSectionName, filepath)
	}

	// Fetch the MFA device ID entry - if there is one
	key := defaultSection.Key(mfaDeviceIDKey)
	if key == nil {
		return "", fmt.Errorf("%s key not found in default section of %s", mfaDeviceIDKey, filepath)
	}

	// Return the value of the key
	return key.String(), nil
}
