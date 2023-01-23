package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewGenerateCmd())
}

var rootCmd = &cobra.Command{
	Use:     "ptg",
	Short:   "ptg is a tool for generating proto file through go",
	Long:    "ptg is a tool for generating proto file through go",
	Version: "0.0.1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
