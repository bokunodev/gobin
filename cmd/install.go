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

var installCmd = &cobra.Command{
	Use:   "install git.tld/org/mod/v2/cmd/pkg@tag",
	Short: "Install a package",
	Args:  cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		pkg := args[0]
		tag := "latest"
		exe := exename(pkg)

		if n := strings.LastIndex(pkg, "@"); n != -1 {
			tag = pkg[n+1:]
			pkg = pkg[:n]
			exe = exename(pkg)
		}

		conf, err := loadconfig()
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("go", append(prepend(conf.BuildFlags, "install"),
			strings.Join([]string{pkg, tag}, "@"))...)

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

		conf.Modules[exe] = module{Pkg: pkg, Tag: tag}
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
	rootCmd.AddCommand(installCmd)
}
