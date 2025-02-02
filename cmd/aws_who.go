package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	awssts "github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/ch0ppy35/sherlock/internal/aws"
	"github.com/ch0ppy35/sherlock/pkg/aws/sts"
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
		aws.WhoAmI(ctx, client)
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
