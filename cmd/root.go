package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/codeshelldev/goplater/pkg/gplt/utils/tmplerrors"
	"github.com/codeshelldev/goplater/utils/pretty"
	"github.com/codeshelldev/goplater/utils/prettyerrors"
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
		tmplErr, ok := errors.AsType[tmplerrors.TemplateError](err)
		
		if ok {
			prettyErr := prettyerrors.Convert(tmplErr)
			printer.Error(tmplErr.CleanContext() + "\n" + prettyErr.ErrorSnippet(pretty.ColorsEnabled(os.Stdout, false)))
		} else {
			printer.Error(err.Error())
		}

		fmt.Println()

		os.Exit(1)
	}
}