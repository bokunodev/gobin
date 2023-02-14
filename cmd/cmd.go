package cmd

import (
	"encoding/json"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"time"
)

/*
	{
		"Path": "github.com/urfave/cli/v2",
		"Version": "v2.23.6",
		"Time": "2022-12-02T14:14:14Z",
		"Update": {
			"Path": "github.com/urfave/cli/v2",
			"Version": "v2.24.3",
			"Time": "2023-02-01T14:23:16Z"
		},
		"GoMod": "/home/boku/.go/pkg/mod/cache/download/github.com/urfave/cli/v2/@v/v2.23.6.mod",
		"GoVersion": "1.18",
		"Origin": {
		"VCS": "git",
		"URL": "https://github.com/urfave/cli",
		"Ref": "refs/tags/v2.23.6",
		"Hash": "f9652e31767f6bbddb654468654fd42473a9eec0"
		}
	}
*/

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

func pkginfo(pkg string) (module, error) {
	m := module{}

	cmd := exec.Command("go", "list", "-u", "-m", "-json", pkg)
	p, err := cmd.CombinedOutput()
	if err != nil {
		return m, err
	}

	return m, json.Unmarshal(p, &m)
}
