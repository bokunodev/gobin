package cmd

import (
	"encoding/json"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type config struct {
	Modules    map[string]module `json:"modules"`
	BuildFlags []string          `json:"build_flags"`
}

type module struct {
	Pkg string `json:"pkg"`
	Tag string `json:"tag"`
}

var gobinConfigFile = filepath.Join(build.Default.GOPATH, ".gobin.json")

func loadconfig() (config, error) {
	conf := config{
		Modules:    map[string]module{},
		BuildFlags: []string{},
	}

	p, err := os.ReadFile(gobinConfigFile)
	if err != nil {
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

func prepend[T any](dst []T, item T) []T {
	return append(append([]T(nil), item), dst...)
}

func gobin() string {
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		gobin = filepath.Join(build.Default.GOPATH, "bin")
	}

	return gobin
}
