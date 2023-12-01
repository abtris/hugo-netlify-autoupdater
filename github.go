package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/hashicorp/go-version"
)

func isNewVersion(hugoVersion string, netlifyConfigVersion string) bool {
	v1, err := version.NewVersion(hugoVersion)
	if err != nil {
		log.Printf("Error parsing version \"%s\"", hugoVersion)
		return false
	}
	v2, err := version.NewVersion(netlifyConfigVersion)
	if err != nil {
		log.Printf("Error parsing version \"%s\"", netlifyConfigVersion)
		return false
	}

	if v2.LessThan(v1) {
		return true
	}
	return false
}

func getCurrentHugoVersion(ctx context.Context, client *github.Client, owner, repo string) (string, string, string, error) {
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		log.Printf("Get latest release error %v", err)
		return "", "", "", err
	}

	return strings.TrimPrefix(release.GetTagName(), "v"), release.GetHTMLURL(), release.GetBody(), nil
}

func getCurrentDeployedFile(ctx context.Context, client *github.Client,
	owner, repo, path, branch string) (string, error) {

	file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path,
		&github.RepositoryContentGetOptions{Ref: branch},
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

func getCurrentDeployedVersion(ctx context.Context, client *github.Client,
	owner, repo, path, branch string) (string, string, error) {

	content, err := getCurrentDeployedFile(ctx, client, owner, repo, path, branch)
	if err != nil {
		return "", "", err
	}
	config, err := parseNetlifyConf(content)
	if err != nil {
		return "", "", err
	}
	if len(config.Build.BuildEnv.HugoVersion) > 0 {
		return config.Build.BuildEnv.HugoVersion, content, nil
	}
	return "", "", fmt.Errorf("ERROR: Empty version")
}

func getCommitBranch(hugoVersion string) string {
	return fmt.Sprintf("updater/version-%s", hugoVersion)
}

func getRef(ctx context.Context, client *github.Client,
	owner, repo, branch, commitBranch string) (ref *github.Reference,
	isNewBranch bool, err error) {

	var baseRef *github.Reference
	// if branch exists get back ref
	if ref, _, err = client.Git.GetRef(ctx, owner, repo, "refs/heads/"+commitBranch); err == nil {
		return ref, false, nil
	}
	// get base ref (master only supported now)
	if baseRef, _, err = client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch); err != nil {
		return nil, false, err
	}
	// create new branch
	newRef := &github.Reference{
		Ref:    github.String("refs/heads/" + commitBranch),
		Object: &github.GitObject{SHA: baseRef.Object.SHA},
	}
	ref, _, err = client.Git.CreateRef(ctx, owner, repo, newRef)
	return ref, true, err
}

func getTree(ctx context.Context, client *github.Client, owner, repo string,
	ref *github.Reference, filename, source string) (tree *github.Tree, err error) {

	entries := []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{Path: github.String(filename),
		Type:    github.String("blob"),
		Content: github.String(source),
		Mode:    github.String("100644"),
	})
	tree, _, err = client.Git.CreateTree(ctx, owner, repo, *ref.Object.SHA, entries)
	return tree, err
}

func pushCommit(ctx context.Context, client *github.Client, owner, repo string,
	ref *github.Reference, tree *github.Tree, hugoVersion string) (err error) {

	parent, _, err := client.Repositories.GetCommit(ctx, owner, repo, *ref.Object.SHA, &github.ListOptions{})
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA
	commitMessage := fmt.Sprintf("fix(deps): Update Hugo to version %s", hugoVersion)
	commiterName := "Updater-bot"
	commiterEmail := "updater-bot@github.com"
	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{
		Date:  &github.Timestamp{date},
		Name:  &commiterName,
		Email: &commiterEmail,
	}
	commit := &github.Commit{
		Author:  author,
		Message: &commitMessage,
		Tree:    tree,
		Parents: []*github.Commit{parent.Commit},
	}
	opts := &github.CreateCommitOptions{}
	newCommit, _, err := client.Git.CreateCommit(ctx, owner, repo, commit, opts)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, owner, repo, ref, false)
	return err
}

func createPullRequest(ctx context.Context, client *github.Client,
	owner, repo, branch, hugoVersion, releaseURL, releaseInfo, commitBranch string) error {

	prBranch := commitBranch
	prSubject := fmt.Sprintf("[hugo-updater] Update Hugo to version %s", hugoVersion)
	prDescription := fmt.Sprintf("%s\nMore details in %s\n\n%s", prSubject, releaseURL, releaseInfo)
	baseBranch := branch
	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &prBranch,
		Base:                &baseBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())

	return nil
}

func getRepoPath(path string) (owner, repo string) {
	paths := strings.Split(path, "/")

	return paths[0], paths[1]
}
