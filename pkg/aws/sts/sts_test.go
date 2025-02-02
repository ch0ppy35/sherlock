package sts

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/assert"
)

var getCallerIdentityTestOutput = &sts.GetCallerIdentityOutput{
	Account: aws.String("123456789012"),
	Arn:     aws.String("arn:aws:iam::123456789012:user/test-user"),
	UserId:  aws.String("AID123EXAMPLE"),
}

func Test_getCallerIdentity(t *testing.T) {
	type args struct {
		ctx    context.Context
		client STSAPI
	}
	tests := []struct {
		name    string
		args    args
		want    *sts.GetCallerIdentityOutput
		wantErr bool
	}{
		{
			name: "Successful call to GetCallerIdentity",
			args: args{
				ctx: context.Background(),
				client: &MockSTSClient{
					GetCallerIdentityFunc: func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
						return getCallerIdentityTestOutput, nil
					},
				},
			},
			want:    getCallerIdentityTestOutput,
			wantErr: false,
		},
		{
			name: "Error in GetCallerIdentity",
			args: args{
				ctx: context.Background(),
				client: &MockSTSClient{
					GetCallerIdentityFunc: func(ctx context.Context, input *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error) {
						return nil, errors.New("Access Denied")
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCallerIdentity(tt.args.ctx, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCallerIdentity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
