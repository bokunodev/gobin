/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed packages",
	Args:  cobra.ExactArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		conf, err := loadconfig()
		if err != nil {
			log.Fatal(err)
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', tabwriter.DiscardEmptyColumns)

		fmt.Fprintf(tw, "package\t|\tpath\t|\tversion\t|\tupdate\n")
		fmt.Fprintf(tw, "----------\t|\t----------\t|\t----------\t|\t----------\n")
		for pkg, mod := range conf.Modules {
			newerVersion := "-"
			if mod.Update != nil {
				newerVersion = mod.Update.Version
			}

			if _, err = fmt.Fprintf(tw, "%s\t|\t%s\t|\t%s\t|\t%s\n", pkg, mod.Path, mod.Version, newerVersion); err != nil {
				log.Fatal(err)
			}
		}

		if err = tw.Flush(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
