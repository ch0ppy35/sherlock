package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "N/A"
	date    = "NOW"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Sherlock's version",
	Long:  "The version command will display Sherlock's version",
	Run: func(cmd *cobra.Command, args []string) {
		printVersionInfo()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersionInfo() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter

	t.SetTitle("Sherlock Info")
	t.AppendRows([]table.Row{
		{"App Version", version},
		{"Commit", commit},
		{"Build Time", date},
	})
	t.Render()
}
