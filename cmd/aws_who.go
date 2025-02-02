package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	awssts "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/ch0ppy35/sherlock/pkg/sts"
	"github.com/spf13/cobra"
)

// whoCmd represents the dns command
var whoCmd = &cobra.Command{
	Use:                   "who",
	DisableFlagsInUseLine: true,
	Short:                 "Show AWS identity details",
	Long:                  "Retrieve the current AWS caller identity, including account, ARN, and user ID.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, client := setup()
		whoAmI(ctx, client)
	},
}

func init() {
	awsCmd.AddCommand(whoCmd)
}

func setup() (context.Context, *sts.STSClient) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Printf("Unable to load AWS SDK config: %v", err)
		os.Exit(1)
	}

	return ctx, &sts.STSClient{Client: awssts.NewFromConfig(cfg)}
}

func whoAmI(ctx context.Context, client *sts.STSClient) {
	who, err := client.GetCallerIdentity(ctx, &awssts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Printf("something went wrong, are you sure you're using an active session? err: %v", err)
		os.Exit(1)
	}

	profile, exists := os.LookupEnv("AWS_PROFILE")
	if !exists {
		fmt.Println("AWS_PROFILE is not set. Default profile in use")
		profile = "default"
	}

	fmt.Printf("AWS Profile: %s\n", profile)
	fmt.Printf("Account: %s\n", *who.Account)
	fmt.Printf("Arn: %s\n", *who.Arn)
	fmt.Printf("UserId: %s\n", *who.UserId)
}
