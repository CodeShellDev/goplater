package cmd

import (
	"os"

	"github.com/spf13/cobra"
)


var verbose bool
var supress bool

var rootCmd = &cobra.Command{
	Use:   "goplater",
    Short: "Go Template CLI",
    Long:  `Go CLI Programm to Template files.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print additional information")
	rootCmd.Flags().BoolVarP(&supress, "ignore-errors", "i", false, "ignore / supress errors")
}