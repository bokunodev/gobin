/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var reinstallCmd = &cobra.Command{
	Use:   "reinstall exe",
	Short: "Reinstall a package",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		conf, err := loadconfig()
		if err != nil {
			log.Fatal(err)
		}

		mod, ok := conf.Modules[args[0]]
		if !ok {
			log.Fatalf("%q was not installed with gobin", args[0])
		}

		cmd := exec.Command("go", append(prepend(conf.BuildFlags, "install"),
			strings.Join([]string{mod.Pkg, mod.Tag}, "@"))...)

		if ok, _ := c.Flags().GetBool("debug"); ok {
			log.Println(cmd.String())
		}

		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if code := cmd.ProcessState.ExitCode(); code != 0 {
			log.Fatal(cmd.ProcessState.String(), code)
		}

		conf.Modules[args[0]] = mod
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
	rootCmd.AddCommand(reinstallCmd)
}
