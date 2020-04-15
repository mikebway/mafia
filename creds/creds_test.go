package creds

// Copyright © 2020 Michael D Broadway <mikebway@mikebway.com>
//
// Licensed under the ISC License (ISC)
//
// See creds.go for overall package documentation. This file contains
// unit tests for the creds.go functions.

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aws/aws-sdk-go/service/sts"
)

// TestGetSessionCredentialsSuccess substitutes a mock wrapper function for the
// AWS STS GetSessionToken(..) call so that we can guarantee success and see what
// happens.
func TestGetSessionCredentialsSuccess(t *testing.T) {

	// Put the package back into its normal state after we are done with the test
	defer ResetPackageDefaults()

	// Set up a mock AWS STS wrapper function
	accessKey := "key"
	secret := "secret"
	token := "token"
	expiration := time.Now()
	SetGetSessionTokenFunc(func(awsService *sts.STS, input *sts.GetSessionTokenInput) (*sts.GetSessionTokenOutput, error) {
		return &sts.GetSessionTokenOutput{
				Credentials: &sts.Credentials{
					AccessKeyId:     &accessKey,
					SecretAccessKey: &secret,
					SessionToken:    &token,
					Expiration:      &expiration,
				},
			},
			nil
	})

	// Invoke our test target
	credentials, err := GetSessionCredentials("mfa-device-id", "123456", 3600)
	require.Nil(t, err, "there should have been no error")
	require.Equal(t, accessKey, *credentials.AccessKeyID, "Access key did not match expected value")
	require.Equal(t, secret, *credentials.SecretAccessKey, "Secret did not match expected value")
	require.Equal(t, token, *credentials.SessionToken, "session token did not match expected value")
}

// TestGetSessionCredentialsFailure invokes GetSessionCredentials(..) without mocking
// the AWS STS GetSessionToken(..) call wrapper to test what happens when AWS is really
// called under circoumstances where we can be certain that the request will be rejected.
func TestGetSessionCredentialsFailure(t *testing.T) {

	// Invoke our test target with an utterly bogus MFA device serial number and token
	credentials, err := GetSessionCredentials("mfa-device-id", "123456", 3600)
	require.NotNil(t, err, "there should have an error")
	require.Nil(t, credentials, "no credentials should have been obtained")
}
