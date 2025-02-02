package cmd

import (
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
		ctx, client := sts.ClientConnect()
		aws.WhoAmI(ctx, client)
	},
}

func init() {
	awsCmd.AddCommand(whoCmd)
}
