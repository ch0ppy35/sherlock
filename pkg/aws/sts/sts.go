package sts

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type GetCallerIdentityInput = sts.GetCallerIdentityInput
type GetCallerIdentityOutput = sts.GetCallerIdentityOutput

// STS Client
type STSAPI interface {
	GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

type STSClient struct {
	Client *sts.Client
}

func (c *STSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return c.Client.GetCallerIdentity(ctx, input, optFns...)
}

func ClientConnect() (context.Context, *STSClient) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Unable to load AWS SDK config: %v", err)
		os.Exit(1)
	}

	return ctx, &STSClient{Client: sts.NewFromConfig(cfg)}
}

// STS Client Mocks
type MockSTSClient struct {
	GetCallerIdentityFunc func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

func (m *MockSTSClient) GetCallerIdentity(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
	return m.GetCallerIdentityFunc(ctx, input, optFns...)
}
