package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

var client *github.Client

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	// getConfig
	conf, err := parseConfigFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Source repo: %s", conf.SourceRepoReleases)
	// github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	// getCurrentHugoVersion
	hugoVersion, err := getCurrentHugoVersion(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hugoVersion)
	// for _, repository := range conf.TargetRepository {}
	// getCurrentDeployedVersion for all config.targetRepos (done)
	// compareVersion (done)
	// preparePR (getRef, done)
	// createPR (done)
}
