/*
Copyright © 2023 bokunodev bokunocode@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade exe",
	Short: "Upgrade installed package",
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

		curMod, err := pkginfo(strings.Join([]string{mod.Path, mod.Version}, "@"))
		if err != nil {
			log.Fatal(err)
		}

		if curMod.Update == nil {
			fmt.Printf("no newer version available for %q\n", args[0])
			return
		}

		goargs := []string{"install"}
		goargs = append(goargs, conf.BuildFlags...)
		goargs = append(goargs,
			strings.Join([]string{
				curMod.Update.Path,
				curMod.Update.Version,
			}, "@"))

		cmd := exec.Command("go", goargs...)

		if ok, _ := c.InheritedFlags().GetBool("debug"); ok {
			log.Println(cmd.String())
		}

		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if code := cmd.ProcessState.ExitCode(); code != 0 {
			log.Fatal(cmd.ProcessState.String(), code)
		}

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
	rootCmd.AddCommand(upgradeCmd)
}
