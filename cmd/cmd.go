package cmd

import (
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type module struct {
	Path      string    `json:"Path,omitempty"`
	Version   string    `json:"Version,omitempty"`
	Time      time.Time `json:"Time,omitempty"`
	Update    *module   `json:"Update,omitempty"`
	GoMod     string    `json:"GoMod,omitempty"`
	GoVersion string    `json:"GoVersion,omitempty"`
	Origin    origin    `json:"Origin,omitempty"`
}

type origin struct {
	Vcs  string `json:"VCS,omitempty"`
	URL  string `json:"URL,omitempty"`
	Ref  string `json:"Ref,omitempty"`
	Hash string `json:"Hash,omitempty"`
}

type config struct {
	Modules    map[string]module `json:"modules"`
	BuildFlags []string          `json:"build_flags"`
}

var gobinConfigFile = filepath.Join(build.Default.GOPATH, ".gobin.json")

func loadconfig() (config, error) {
	conf := config{
		Modules:    map[string]module{},
		BuildFlags: []string{},
	}

	p, err := os.ReadFile(gobinConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return conf, err
	}

	return conf, json.Unmarshal(p, &conf)
}

var hasSemverSuffix = regexp.MustCompile(`/v\w+$`)

func exename(s string) string {
	if hasSemverSuffix.MatchString(s) {
		return path.Base(hasSemverSuffix.ReplaceAllString(s, ""))
	}

	return path.Base(s)
}

func gobin() string {
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		gobin = filepath.Join(build.Default.GOPATH, "bin")
	}

	return gobin
}

func pkginfo(pkg, tag string) (module, error) {
	m := module{}

	pkgpath := pkg

begin:
	cmd := exec.Command("go", "list", "-u", "-m", "-json",
		strings.Join([]string{pkgpath, tag}, "@"))
	p, err := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() != 0 {
		pkgpath = path.Dir(pkgpath)
		if pkgpath == "." {
			return m, fmt.Errorf(`not found: module %s: no matching versions for query "%s"`, pkg, tag)
		}
		goto begin
	}
	if err != nil {
		return m, err
	}

	return m, json.Unmarshal(p, &m)
}
