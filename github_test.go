package main

import (
	"context"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/hashicorp/go-version"
)

func TestGetCurrentHugoVersion(t *testing.T) {
	t.Parallel()
	expected := "0.81.0"

	var client *github.Client
	var ctx = context.Background()
	// for public repo, you don't need credentials
	client = github.NewClient(nil)
	// public repo as source
	sourceOwner := "gohugoio"
	sourceRepo := "hugo"
	real, _, err := getCurrentHugoVersion(ctx, client, sourceOwner, sourceRepo)
	if err != nil {
		t.Fatalf("Get error %v", err)
	}
	expectedVersion, _ := version.NewVersion(expected)
	realVersion, _ := version.NewVersion(real)
	if expectedVersion.GreaterThanOrEqual(realVersion) {
		t.Errorf("Real version %v is greater than expected %v)", realVersion, expectedVersion)
	}
}

func TestGetCurrentDeployedVersion(t *testing.T) {
	t.Parallel()
	expected := "0.83.1"

	var client *github.Client
	var ctx = context.Background()

	client = github.NewClient(nil)

	real, _, err := getCurrentDeployedVersion(ctx, client, "abtris", "12ApiaryTest", "netlify.toml", "master")
	if err != nil {
		t.Fatalf("Get error %v", err)
	}

	if real != expected {
		t.Errorf("Expected %v and real %v)", expected, real)
	}
}

func TestIsNewVersion(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		hugoVersion    string
		netlifyVersion string
		result         bool
	}{
		{name: "Equal", hugoVersion: "0.10.1", netlifyVersion: "0.10.1", result: false},
		{name: "Lower", hugoVersion: "0.10.1", netlifyVersion: "0.10.2", result: false},
		{name: "New", hugoVersion: "0.10.2", netlifyVersion: "0.10.1", result: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isNewVersion(test.hugoVersion, test.netlifyVersion)
			if result != test.result {
				t.Errorf("Expected %v and real %v)", test.result, result)
			}
		})
	}
}

func TestGetRepoPath(t *testing.T) {
	t.Parallel()
	input := "owner/repo"
	expectedOwner := "owner"
	expectedRepo := "repo"
	owner, repo := getRepoPath(input)

	if owner != expectedOwner {
		t.Errorf("Expected %v and real %v)", expectedOwner, owner)
	}
	if repo != expectedRepo {
		t.Errorf("Expected %v and real %v)", expectedRepo, repo)
	}
}
