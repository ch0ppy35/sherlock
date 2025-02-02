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

func GetCallerIdentity(ctx context.Context, client STSAPI) (*sts.GetCallerIdentityOutput, error) {
	return client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
}

func (c *STSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return c.Client.GetCallerIdentity(ctx, input, optFns...)
}

// Mocks
type MockSTSClient struct {
	GetCallerIdentityFunc func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

func (m *MockSTSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return m.GetCallerIdentityFunc(ctx, input, optFns...)
}
