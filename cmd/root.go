package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/codeshelldev/goplater/utils/pretty"
	"github.com/codeshelldev/goplater/utils/printer"
	"github.com/codeshelldev/goplater/utils/tmplerrors"
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
		tmplErr, ok := errors.AsType[tmplerrors.TemplateError](err)
		
		if ok {
			printer.Error(tmplErr.CleanContext() + "\n" + tmplErr.ErrorSnippet(pretty.ColorsEnabled(os.Stdout, false)))
		} else {
			printer.Error(err.Error())
		}

		fmt.Println()

		os.Exit(1)
	}
}