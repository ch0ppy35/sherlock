package sts

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STSAPI interface {
	GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type STSClient struct {
	Client *sts.Client
}

func (c *STSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return c.Client.GetCallerIdentity(ctx, input, optFns...)
}

// Mocks
type MockSTSClient struct {
	getCallerIdentityFunc func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

func (m *MockSTSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return m.getCallerIdentityFunc(ctx, input, optFns...)
}

//
//
