package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/ch0ppy35/sherlock/pkg/aws/sts"
	"github.com/ch0ppy35/sherlock/pkg/ui"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// getAWSProfile retrieves the AWS_PROFILE environment variable
func getAWSProfile() string {
	profile, exists := os.LookupEnv("AWS_PROFILE")
	if !exists {
		fmt.Fprintf(
			os.Stdout,
			"%s",
			ui.DefaultColorWriters["magenta"]("AWS_PROFILE is not set, default profile in use\n"),
		)
		return "default"
	}
	return profile
}

// renderWhoAmITable prints the AWS identity information in a table format
func renderWhoAmITable(profile string, who *sts.GetCallerIdentityOutput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter

	t.SetTitle("AWS Session Info")
	t.AppendRows([]table.Row{
		{"AWS Profile", profile},
		{"Account", *who.Account},
		{"Arn", *who.Arn},
		{"UserId", *who.UserId},
	})

	t.Render()
}

// WhoAmI retrieves and prints AWS identity details
func WhoAmI(ctx context.Context, client sts.STSAPI) {
	profile := getAWSProfile()

	who, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		ui.PrintErrMsgWithStatus("ERROR", "red", "%v", err)
		os.Exit(1)
	}

	renderWhoAmITable(profile, who)
}
