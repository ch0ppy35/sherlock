package sts

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
)

var getCallerIdentityTestOutputGood = &sts.GetCallerIdentityOutput{
	Account: aws.String("123456789012"),
	Arn:     aws.String("arn:aws:iam::123456789012:user/test-user"),
	UserId:  aws.String("AID123EXAMPLE"),
}

func Test_getCallerIdentity(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		client  STSAPI
		want    *sts.GetCallerIdentityOutput
		wantErr bool
	}{
		{
			name: "Successful call to GetCallerIdentity",
			ctx:  context.Background(),
			client: &MockSTSClient{
				GetCallerIdentityFunc: func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
					return getCallerIdentityTestOutputGood, nil
				},
			},
			want:    getCallerIdentityTestOutputGood,
			wantErr: false,
		},
		{
			name: "Error in GetCallerIdentity",
			ctx:  context.Background(),
			client: &MockSTSClient{
				GetCallerIdentityFunc: func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
					return nil, errors.New("operation error STS: GetCallerIdentity, get identity: get credentials: failed to refresh cached credentials, no EC2 IMDS role found, operation error ec2imds: GetMetadata, request canceled, context deadline exceededexit status 1")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.GetCallerIdentity(tt.ctx, &sts.GetCallerIdentityInput{})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCallerIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
