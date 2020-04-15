// Package creds interfaces with AWS to obtain new session credentials.
//
// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
package creds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// SessionCredentials wraps the AWS credentials obtained for and API authentication session
type SessionCredentials struct {
	AccessKeyID     *string // The access key ID that identifies the temporary security credentials
	SecretAccessKey *string
	SessionToken    *string
}

// GetSessionTokenFunc is a function type that corresponds to the AWS STS function for obtaining
// a session token. Rather than calling this function directly from GetSessionCredentials(..),
// it is called via a function variable; when unit testing, this function variable can be
// be overriden and point to a mock implementation.
type GetSessionTokenFunc func(awsService *sts.STS, input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error)

var (

	// A function variable that, normally, wraps the AWS STS GetSessionToken(..) function
	// but can be overriden for unit testsing. This is initialied at load time via a call to
	// the ResetPackageDefaults(..) function.
	getSessionTokenFunc GetSessionTokenFunc
)

// Load time initialization
func init() {

	// Configure the default state of this package, in particular, set the getSessionTokenFunc
	// to reference a function that wraps the AWS STS GetSessionToken(..) function.
	ResetPackageDefaults()
}

// GetSessionCredentials combines AWS credentials from the environment with a provided MFA
// token to authenticate with AWS and obtain session credentials.
//
// The mfaSerialNumber number identifies the device that was used to source the
// mfaToken, the latter typically being a 6 digit number while the device serial
// number really is a string in the form arn:aws:iam::accountnumber:mfa/username.
//
// Duration is a time period expressed as a number of seconds. AWS accepts values
// between 900 seconds (15 minutes) to 129,600 seconds (36 hours).
func GetSessionCredentials(mfaSerialNumber, mfaToken string, duration int64) (*SessionCredentials, error) {

	// Obtain an AWS STS client
	svc := sts.New(session.New())

	// Prep the input structure for the get session request
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(duration),
		SerialNumber:    aws.String(mfaSerialNumber),
		TokenCode:       aws.String(mfaToken),
	}

	// Request a new session from AWS via our wrapper function variable.
	// When unit testing, the
	result, err := getSessionTokenFunc(svc, input)
	if err != nil {
		return nil, err
	}

	// Translate the result into our own format that does not require the caller
	// to also import the AWS STS package and return that
	return &SessionCredentials{
		AccessKeyID:     result.Credentials.AccessKeyId,
		SecretAccessKey: result.Credentials.SecretAccessKey,
		SessionToken:    result.Credentials.SessionToken,
	}, nil
}

// SetGetSessionTokenFunc allows unit tests to substitute a mock function in place of
// the default AWS STS GetSessionToken(..) wrapper so that tests can control the responses.
func SetGetSessionTokenFunc(f GetSessionTokenFunc) {
	getSessionTokenFunc = f
}

// ResetPackageDefaults establishes or reestablishes the normal package global values.
// This is called during package initialization and alos by unit tests needing to
// leave the package as they found it.
func ResetPackageDefaults() {

	// Configure the function wrapper used to ask AWS STS for a session token
	getSessionTokenFunc = func(awsService *sts.STS, input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
		return awsService.GetSessionToken(input)
	}
}
