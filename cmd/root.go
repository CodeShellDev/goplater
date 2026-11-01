package cmd

import (
	"os"

	"github.com/codeshelldev/goplater/utils/printer"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goplater",
    Short: "Go Template CLI",
    Long:  "Go CLI programm for templating files.",
}

func Execute() {
	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	
	err := rootCmd.Execute()
	if err != nil {
		printer.Error(printer.FormatError(err))
		os.Exit(1)
	}
}