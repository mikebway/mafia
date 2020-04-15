package mfile

// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
//
// See doc.go for other overall package documentation. This file contains
// unit tests for the read.go functions as well as the common test setup and
// tear down functions used by most package tests.

import (
	"fmt"
	"os"
	"testing"

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

// TestGetMFADeviceID examines the happy path where an MFA device serial number has been
// stored in the AWS credentials and can be retrieved successfuly.
func TestGetMFADeviceID(t *testing.T) {

	// Revert the package state back to notmal after the test has run
	defer ResetPackageDefaults()

	// Establisd a fake credentials file with an MFA device ID
	setFakeCredentials(DefaultSectionName, fakeMFADeviceID)

	// Asking for the MFA device ID should fail
	id, err := GetMFADeviceID()
	require.Nil(t, err, "there should not have been an error")
	require.Equal(t, fakeMFADeviceID, id, "not the expected MFA device ID")
}

// TestGetMFADeviceIDMissingKey examines the sad path where an MFA device serial number has not
// been stored in the AWS credentials and can be retrieved successfuly.
func TestGetMFADeviceIDMissingKey(t *testing.T) {

	// Revert the package state back to notmal after the test has run
	defer ResetPackageDefaults()

	// Establisd a fake credentials file with no MFA device ID
	setFakeCredentials(DefaultSectionName, "")

	// Asking for the MFA device ID should fail
	id, err := GetMFADeviceID()
	require.NotNil(t, err, "there should have been an error")
	require.Equal(t, "mfa_device_id key not found in default section of ./credentials.test", err.Error(), "not the expected error")
	require.Empty(t, id, "no MFA device ID should have been returned")
}

// TestGetMFADeviceIDMissingSection examines the sad path where the default section has not
// been stored in the AWS credentials and can be retrieved successfuly.
func TestGetMFADeviceIDMissingSection(t *testing.T) {

	// Revert the package state back to notmal after the test has run
	defer ResetPackageDefaults()

	// Establisd a fake credentials file with no default section
	setFakeCredentials("not-the-droids", "")

	// Asking for the MFA device ID should fail
	id, err := GetMFADeviceID()
	require.NotNil(t, err, "there should have been an error")
	require.Equal(t, "default section not found in ./credentials.test", err.Error(), "not the expected error")
	require.Empty(t, id, "no MFA device ID should have been returned")
}

// TestGetMFADeviceIDMissingFile examines the sad path where the expected AWS credentials
// file does not exist
func TestGetMFADeviceIDMissingFile(t *testing.T) {

	// Revert the package state back to notmal after the test has run
	defer ResetPackageDefaults()

	// Have the package think the credentials file is in a place where there are no files
	OverrideDefaultCredentialsFilepath("/you/got/no/skin/on/me-cos-i-do-not-exist")

	// Asking for the MFA device ID should fail
	id, err := GetMFADeviceID()
	require.NotNil(t, err, "there should have been an error")
	require.Contains(t, err.Error(), "no such file or directory", "not the expected error")
	require.Empty(t, id, "no MFA device ID should have been returned")
}

// setFakeCredentials populates a fake AWS credentials file in the current
// working directory, with or without an MFA device serial number / ID. The
// package globals are then manipulated such that this fake file will be used
// any future test execution.
func setFakeCredentials(sectionName, mfaDeviceID string) {

	// Start with an empty configuration file content structure
	cfg := ini.Empty()

	// Populate it with the basics
	defaultSection, err := cfg.NewSection(sectionName)
	defaultSection.NewKey(AccessKeyIDKey, fakeAccessKeyID)
	defaultSection.NewKey(SecretAccessKeyKey, fakeSecretAccessKey)

	// If we were given an MFA device ID, put that in our configuration too
	if len(mfaDeviceID) != 0 {
		defaultSection.NewKey(MfaDeviceIDKey, mfaDeviceID)
	}

	// Write the file
	err = cfg.SaveTo(fakeCredentialsFilePath)

	// That really should not faile to wrote, but if it did abort the tests
	// cos nothing will work after this
	if err != nil {
		fmt.Printf("\nFailed to write fake credentials: %s\n", err.Error())
		os.Exit(999)
	}

	// All looks good - trick the package into using the fake file we just wrote
	OverrideDefaultCredentialsFilepath(fakeCredentialsFilePath)
}
