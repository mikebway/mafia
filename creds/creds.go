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

	// Request a new session from AWS
	result, err := svc.GetSessionToken(input)
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
