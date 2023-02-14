/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gobin",
	Short: "Help you manage tools installed trough `go install` command",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().Bool("debug", false, "Log generated `go install` command")
}
