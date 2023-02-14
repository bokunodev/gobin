/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove exe",
	Short: "Remove installed package",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		binpath := gobin()

		conf, err := loadconfig()
		if err != nil {
			log.Fatal(err)
		}

		if _, ok := conf.Modules[args[0]]; !ok {
			log.Fatalf("%q was not installed with gobin", args[0])
		}

		if err = os.Remove(filepath.Join(binpath, args[0])); err != nil {
			log.Fatal(err)
		}

		delete(conf.Modules, args[0])

		p, err := json.MarshalIndent(conf, "", "\t")
		if err != nil {
			log.Fatal(err)
		}

		if err = os.WriteFile(gobinConfigFile, p, 0o644); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
