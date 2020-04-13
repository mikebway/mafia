package mfile

import (
	"fmt"

	"github.com/go-ini/ini"
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
	sessionSection := cfg.Section(sessionSectionName)

	// Set the section key/values
	sessionSection.NewKey(accessKeyIDKey, *accessKeyID)
	sessionSection.NewKey(secretAccessKeyKey, *secretAccessKey)
	sessionSection.NewKey(sessionTokenKey, *sessionToken)

	// Save the file and we are done
	return cfg.SaveTo(filepath)
}
