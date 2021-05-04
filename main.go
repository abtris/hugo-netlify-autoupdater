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
	hugoVersion, releaseUrl, err := getCurrentHugoVersion(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hugoVersion)
	fmt.Println(releaseUrl)

	for _, repository := range conf.TargetRepository {
		fmt.Printf("%s\n%s\n%s\n", repository.Repo, repository.TargetFile, repository.TargetVariable)
		// getCurrentDeployedVersion for all config.targetRepos (done)
		owner, repo := getRepoPath(repository.Repo)
		deployVersion, deployContent, err := getCurrentDeployedVersion(ctx, client, owner, repo, repository.TargetFile)
		if err != nil {
			fmt.Printf("Error in %v", err)
		}
		if isNewVersion(hugoVersion, deployVersion) {
			updatedContent := updateVersion(hugoVersion, deployContent)
			fmt.Println(updatedContent)
			// preparePR (getRef, done)
			// createPR (done)
		} else {
			fmt.Printf("No new version in %s/%s (current: %s)\n", owner, repo, deployVersion)
		}
	}
}
