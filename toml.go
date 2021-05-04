package main

import (
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
)

type repository struct {
	Repo           string `toml:"repo"`
	TargetFile     string `toml:"target_file"`
	TargetVariable string `toml:"target_variable"`
	Branch         string `toml:"branch"`
}

type config struct {
	SourceRepoReleases string       `toml:"source_repo_releases"`
	TargetRepository   []repository `toml:"target_repos"`
}

type netlifyConfig struct {
	Build   netlifyBuild `toml:"build"`
	Context netlifyBuild `toml:"context"`
}

type netlifyBuild struct {
	Command  string                  `toml:"command"`
	BuildEnv netlifyBuildEnvironment `toml:"environment"`
}
type netlifyBuildEnvironment struct {
	HugoVersion string `toml:"HUGO_VERSION"`
}

func parseConfigFile(filepath string) (config, error) {
	var conf config
	if _, err := toml.DecodeFile(filepath, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func parseNetlifyConfFile(filepath string) (netlifyConfig, error) {
	var conf netlifyConfig
	if _, err := toml.DecodeFile(filepath, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func parseNetlifyConf(content string) (netlifyConfig, error) {
	var conf netlifyConfig
	if _, err := toml.Decode(content, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func updateVersion(hugoVersion, deployContent string) string {

	regexp, err := regexp.Compile(`HUGO_VERSION = \"(\d+\.\d+\.\d+)\"`)
	if err != nil {
		fmt.Println(err)
	}
	replacement := fmt.Sprintf("HUGO_VERSION = \"%s\"", hugoVersion)

	return regexp.ReplaceAllString(deployContent, replacement)
}
