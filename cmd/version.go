package cmd

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var (
	arch      = "N/A"
	version   = "dev"
	commit    = "N/A"
	date      = "NOW"
	goversion = "N/A"
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

func printVersionInfo() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleRounded)
	t.Style().Title.Align = text.AlignCenter

	t.SetTitle("Sherlock Info")
	t.AppendRows([]table.Row{
		{"Version", version},
		{"Revision", commit},
		{"Arch", arch},
		{"BuildTime", date},
		{"BuildGoVersion", goversion},
	})
	t.Render()
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
