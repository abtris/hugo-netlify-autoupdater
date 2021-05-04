package main

import (
	"log"
	"os"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	// getConfig
	// github client
	// ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	// tc := oauth2.NewClient(ctx, ts)
	// client = github.NewClient(tc)
	// getCurrentHugoVersion (done)
	// getCurrentDeployedVersion for all config.targetRepos (done)
	// compareVersion (done)
	// preparePR
	// checkPRIfExists (don't create every day)
	// createPR
}
