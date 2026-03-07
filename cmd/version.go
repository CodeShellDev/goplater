package cmd

import (
	"github.com/codeshelldev/goplater/internals/store"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
    Short: "Output version",
    Long:  `Output version of binary build.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(store.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}