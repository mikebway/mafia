package mfile

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
//
// See doc.go for other overall package documentation. This file contains
// package methods related to updating the AWS credentials file.

// SaveSessionCredentials writes the given credentials to a "session" section of the
// default AWS credentials file, i.e. $HOME/.aws/credentials.
func SaveSessionCredentials(accessKeyID, secretAccessKey, sessionToken *string) error {

	// Have our siblings do all the work!
	return SaveSessionCredentialsToFile(defaultCredentialsFilePath,
		accessKeyID, secretAccessKey, sessionToken)
}

// SaveSessionCredentialsToFile saves the given credentials to a "session" section of the
// the given AWS credentials file.
func SaveSessionCredentialsToFile(filepath string, accessKeyID, secretAccessKey, sessionToken *string) error {

	// Load the current file contents
	cfg, err := ini.Load(filepath)
	if err != nil {
		return fmt.Errorf("Could not read from credentials file %s: %v", filepath, err)
	}

	// Either load any previously existing session or create a new one with the required name
	sessionSection := cfg.Section(SessionSectionName)

	// Set the section key/values
	sessionSection.NewKey(AccessKeyIDKey, *accessKeyID)
	sessionSection.NewKey(SecretAccessKeyKey, *secretAccessKey)
	sessionSection.NewKey(SessionTokenKey, *sessionToken)

	// Save the file and we are done
	return cfg.SaveTo(filepath)
}
