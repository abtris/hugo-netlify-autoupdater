package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v60/github"
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
		log.Fatalf("Missing or wrong config.toml - %v", err)
	}
	log.Printf("Source repo: %s\n", conf.SourceRepoReleases)
	// github client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
	// getCurrentHugoVersion
	sourceOwner, sourceRepo := getRepoPath(conf.SourceRepoReleases)
	hugoVersion, releaseURL, releaseInfo, err := getCurrentHugoVersion(ctx, client, sourceOwner, sourceRepo)
	if err != nil {
		log.Fatal(err)
	}

	for _, repository := range conf.TargetRepository {
		// getCurrentDeployedVersion for all config.targetRepos (done)
		owner, repo := getRepoPath(repository.Repo)
		deployVersion, deployContent, err := getCurrentDeployedVersion(ctx, client, owner, repo, repository.TargetFile, repository.Branch)
		if err != nil {
			log.Fatalf("Error in getCurrentDeployedVersion - %v", err)
		}
		messageHugoVersion := []byte(hugoVersion)
		errHugo := os.WriteFile(".hugo-version", messageHugoVersion, 0644)
		if errHugo != nil {
			log.Fatal(err)
		}
		messageDeployVersion := []byte(deployVersion)
		errDeployed := os.WriteFile(".deployed-version", messageDeployVersion, 0644)
		if errDeployed != nil {
			log.Fatal(err)
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
				errPR := createPullRequest(ctx, client, owner, repo, repository.Branch, hugoVersion, releaseURL, releaseInfo, commitBranch)
				if errPR != nil {
					log.Fatalf("Error in createPullRequest %v", errPR)
				}
			} else {
				log.Printf("PR branch (%s) already exists.\n", commitBranch)
			}
		} else {
			log.Printf("No new version in %s/%s (current: %s)\n", owner, repo, deployVersion)
		}
	}
}
