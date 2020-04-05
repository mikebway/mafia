// Package file loads and saves AWS credential information from a property file.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package file

import (
	"fmt"

	"github.com/go-ini/ini"
)

const {
	DefaultCredentialsPath = "$HOME/.aws/credentials"
	DefaultSectionName = "defualt"
	MFADeviceIDKey = "mfa_device_id"
}

// GetMFADeviceID attempts to find an MFA device ID in the default section of
// the AWS credentials file, returing either the ID or an error.
func GetMFADeviceID() (string, error) {
	return GetMFADeviceIDFromFile(DefaultCredentialsPath)
}

// GetMFADeviceIDFromFile attempts to find an MFA device ID in the default section of
// the given AWS credentials file, returing either the ID or an error.
func GetMFADeviceIDFromFile(filepath string) (string, error) {

	// Load the file
	cfg, err := ini.Load(filepath)
	if err != nil {
		return nil, fmt.Errorf("Could not read from credentials file %s: %v", filepath, err)
	}

	// Fetch the default section - if there is one
	defaultSection := cfg.Section(DefaultSectionName)
	if defaultSection == nil {
		return "", fmt.Errorf("%s section not found in %s", DefaultSectionName, filepath)
	}

	// Fetch the MFA device ID entry - if there is one
	mfaDeviceIDKey := defaultSection.Key(MFADeviceIDKey)
	if mfaDeviceIDKey == nil {
		return nil, fmt.Errorf("%s key not found in default section of %s", MFADeviceIDKey, filepath)
	}

	// Return the value of the key
	return mfaDeviceIDKey.String(), nil
}
