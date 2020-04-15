package mfile

// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
//
// See doc.go for other overall package documentation. This file contains
// unit tests for the write.go functions.

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

// TestSaveSessionCredentials examines the happy path where session credentials are stored
// successfully in an existing file.
func TestSaveSessionCredentials(t *testing.T) {

	// Revert the package state back to normal after the test has run
	defer ResetPackageDefaults()

	// Establisd a virgin fake credentials file with known contents
	setFakeCredentials(DefaultSectionName, fakeMFADeviceID)

	// Write fake credentials to that
	firstAccessKey := "key_1"
	firstSecret := "secret_1"
	firstToken := "token_1"
	err := SaveSessionCredentials(&firstAccessKey, &firstSecret, &firstToken)
	require.Nil(t, err, "there should not have been an error (first save)")

	// Confirm that the session values were written
	verifyConfiguration(t, firstAccessKey, firstSecret, firstToken)

	// That covers what happens when there were no session values in the config before, what about
	// when we try to overwrite them?
	secondAccessKey := "key_1"
	secondSecret := "secret_1"
	secondToken := "token_1"
	err = SaveSessionCredentials(&secondAccessKey, &secondSecret, &secondToken)
	require.Nil(t, err, "there should not have been an error (second save)")

	// Confirm that the session values were written
	verifyConfiguration(t, secondAccessKey, secondSecret, secondToken)
}

// TestSaveToNonExistentFile looks at the sad path where the supposedly pre-existing
// AWS credentials file does not, in fact, exist
func TestSaveToNonExistentFile(t *testing.T) {

	// Revert the package state back to normal after the test has run
	defer ResetPackageDefaults()

	// Have the package think the credentials file is in a place where there are no files
	OverrideDefaultCredentialsFilepath("/you/got/no/skin/on/me-cos-i-do-not-exist")

	// Write fake credentials to that
	firstAccessKey := "key_1"
	firstSecret := "secret_1"
	firstToken := "token_1"
	err := SaveSessionCredentials(&firstAccessKey, &firstSecret, &firstToken)
	require.NotNil(t, err, "saving to a non-existent file should have failed")
}

// verifyConfiguration checks that the test configuration file contains both of the
// expected sections and they they both contain the expected key/values.
func verifyConfiguration(t *testing.T, accessKeyID, secretAccessKey, sessionToken string) {

	// Load the current file contents
	cfg, err := ini.Load(fakeCredentialsFilePath)
	require.Nil(t, err, "error reading the test credentials file")

	// Confirm thet the default credentials are set as expected
	defaultSection, err := cfg.GetSection(DefaultSectionName)
	require.Nil(t, err, "default section not found in credentials file")
	require.Equal(t, defaultSection.Key(AccessKeyIDKey).Value(), fakeAccessKeyID, "default section, unexpected access key value: [%s]", defaultSection.Key(AccessKeyIDKey).Value())
	require.Equal(t, defaultSection.Key(SecretAccessKeyKey).Value(), fakeSecretAccessKey, "default section, unexpected secret key value: [%s]", defaultSection.Key(SecretAccessKeyKey).Value())

	// Confirm that the default-session section is as expected
	sessionSection, err := cfg.GetSection(SessionSectionName)
	require.Equal(t, sessionSection.Key(AccessKeyIDKey).Value(), accessKeyID, "default-session section, unexpected access key value: [%s]", defaultSection.Key(AccessKeyIDKey).Value())
	require.Equal(t, sessionSection.Key(SecretAccessKeyKey).Value(), secretAccessKey, "default-session section, unexpected secret key value: [%s]", defaultSection.Key(SecretAccessKeyKey).Value())
	require.Equal(t, sessionSection.Key(SessionTokenKey).Value(), sessionToken, "default-session section, unexpected session token value: [%s]", defaultSection.Key(SessionTokenKey).Value())
}
