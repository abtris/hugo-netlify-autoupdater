package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v35/github"
)

func getCurrentHugoVersion(ctx context.Context, client *github.Client) (string, error) {
	release, _, err := client.Repositories.GetLatestRelease(ctx, "gohugoio", "hugo")
	if err != nil {
		fmt.Println(err)
	}
	return strings.TrimPrefix(release.GetTagName(), "v"), nil
}

func getCurrentDeployedFile(ctx context.Context, client *github.Client, owner, repo, path string) (string, error) {
	file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path,
		&github.RepositoryContentGetOptions{Ref: "master"},
	)
	if err != nil {
		return "", err
	}
	content, err := file.GetContent()
	if err != nil {
		return "", err
	}
	return content, nil
}

func getCurrentDeployedVersion(ctx context.Context, client *github.Client, owner, repo, path string) (string, error) {
	content, err := getCurrentDeployedFile(ctx, client, owner, repo, path)
	if err != nil {
		return "", err
	}
	config, err := parseNetlifyConf(content)
	if err != nil {
		return "", err
	}
	return config.Build.BuildEnv.HugoVersion, nil
}
