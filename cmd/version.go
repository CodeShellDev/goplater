package cmd

import (
	"github.com/codeshelldev/goplater/internals/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
    Short: "Output build information",
    Long:  `Output information about binary build.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Version: " + buildinfo.Version)
		cmd.Println("Commit: " + buildinfo.Commit)
		cmd.Println("Date: " + buildinfo.Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}