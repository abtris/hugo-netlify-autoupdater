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
	hugoVersion, releaseURL, err := getCurrentHugoVersion(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hugoVersion)
	fmt.Println(releaseURL)

	for _, repository := range conf.TargetRepository {
		// getCurrentDeployedVersion for all config.targetRepos (done)
		owner, repo := getRepoPath(repository.Repo)
		deployVersion, deployContent, err := getCurrentDeployedVersion(ctx, client, owner, repo, repository.TargetFile, repository.Branch)
		if err != nil {
			fmt.Printf("Error in %v", err)
		}
		if isNewVersion(hugoVersion, deployVersion) {
			updatedContent := updateVersion(hugoVersion, deployContent)
			fmt.Println(updatedContent)
			commitBranch := getCommitBranch(hugoVersion)
			ref, newBranch, err := getRef(ctx, client, owner, repo, repository.Branch, commitBranch)
			if err != nil {
				log.Fatalf("Error in getRef %v", err)
			}
			if newBranch {
				tree, err := getTree(ctx, client, owner, repo, ref, "netlify.toml", updatedContent)
				if err != nil {
					log.Fatalf("Error in getTree %v", err)
				}
				errCommit := pushCommit(ctx, client, owner, repo, ref, tree, hugoVersion)
				if errCommit != nil {
					log.Fatalf("Error in pushCommit %v", errCommit)
				}
				errPR := createPullRequest(ctx, client, owner, repo, repository.Branch, hugoVersion, releaseURL, commitBranch)
				if errPR != nil {
					log.Fatalf("Error in createPullRequest %v", errPR)
				}
			} else {
				fmt.Printf("PR branch (%s) already exists.\n", commitBranch)
			}
		} else {
			fmt.Printf("No new version in %s/%s (current: %s)\n", owner, repo, deployVersion)
		}
	}
}
