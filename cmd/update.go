/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update pkg",
	Short: "Check for newer version",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		conf, err := loadconfig()
		if err != nil {
			log.Fatal(err)
		}

		mod, ok := conf.Modules[args[0]]
		if !ok {
			log.Fatalf("%q was not installed with gobin\n", args[0])
		}

		curMod, err := pkginfo(mod.Path, mod.Version)
		if err != nil {
			log.Fatal(err)
		}

		if curMod.Update == nil {
			fmt.Printf("no newer version available for %q\n", args[0])
			return
		}

		fmt.Printf("%s %s => %s\n", curMod.Path, curMod.Version, curMod.Update.Version)

		conf.Modules[args[0]] = curMod
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
	rootCmd.AddCommand(updateCmd)
}
