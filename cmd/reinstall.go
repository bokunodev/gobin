/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var reinstallCmd = &cobra.Command{
	Use:   "reinstall",
	Short: "Rebuild and reisntall a package",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		conf, err := loadconfig()
		if err != nil {
			log.Fatal(err)
		}

		mod, ok := conf.Modules[args[0]]
		if !ok {
			log.Fatalf("%q was not installed with gobin\n", args[0])
		}

		goargs := []string{"install"}
		goargs = append(goargs, conf.BuildFlags...)
		goargs = append(goargs, strings.Join([]string{mod.RealPath, mod.Version}, "@"))

		cmd := exec.Command("go", goargs...)

		if ok, _ := c.InheritedFlags().GetBool("debug"); ok {
			log.Println(cmd.String())
		}

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if code := cmd.ProcessState.ExitCode(); code != 0 {
			log.Fatal(cmd.ProcessState.String(), code)
		}
	},
}

func init() {
	rootCmd.AddCommand(reinstallCmd)
}
