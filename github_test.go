package main

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/v90/github"
	"github.com/hashicorp/go-version"
)

// newTestClient returns a client authenticated with GITHUB_TOKEN when it is
// set, to avoid the 60 requests/hour unauthenticated rate limit. For public
// repos the unauthenticated client works too, just with the lower limit.
func newTestClient(t *testing.T) *github.Client {
	t.Helper()
	var opts []github.ClientOptionsFunc
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		opts = append(opts, github.WithAuthToken(token))
	}
	client, err := github.NewClient(opts...)
	if err != nil {
		t.Fatalf("Get error %v", err)
	}
	return client
}

func TestGetCurrentHugoVersion(t *testing.T) {
	t.Parallel()
	expected := "0.81.0"

	var ctx = context.Background()
	client := newTestClient(t)
	// public repo as source
	sourceOwner := "gohugoio"
	sourceRepo := "hugo"
	real, _, _, err := getCurrentHugoVersion(ctx, client, sourceOwner, sourceRepo)
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

	var ctx = context.Background()
	client := newTestClient(t)

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

func TestSanitizeMentions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "Empty", input: "", expected: ""},
		{name: "StartOfText", input: "@bep fixed it", expected: "bep fixed it"},
		{name: "StartOfLine", input: "notes\n@bep fixed it", expected: "notes\nbep fixed it"},
		{name: "InSentence", input: "Thanks to @jmooring for the fix", expected: "Thanks to jmooring for the fix"},
		{name: "MultipleMentions", input: "@bep @jmooring", expected: "bep jmooring"},
		{name: "MarkdownLink", input: "[@bep](https://github.com/bep)", expected: "[bep](https://github.com/bep)"},
		{name: "InParens", input: "(@bep)", expected: "(bep)"},
		{name: "Emphasized", input: "*@bep*", expected: "*bep*"},
		{name: "Team", input: "cc @gohugoio/maintainers here", expected: "cc gohugoio/maintainers here"},
		{name: "DashedHandle", input: "by @some-user today", expected: "by some-user today"},
		{name: "Email", input: "updater-bot@github.com", expected: "updater-bot@github.com"},
		{name: "URLPath", input: "https://example.com/@bep", expected: "https://example.com/@bep"},
		{name: "NoMention", input: "Update Hugo to version 0.150.0", expected: "Update Hugo to version 0.150.0"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := sanitizeMentions(test.input)
			if result != test.expected {
				t.Errorf("Expected %q and real %q)", test.expected, result)
			}
		})
	}
}

func TestGetPullRequestBody(t *testing.T) {
	t.Parallel()
	expected := "[hugo-updater] Update Hugo to version 0.150.0\nMore details in https://github.com/gohugoio/hugo/releases/tag/v0.150.0\n\nThanks to bep for the fix"
	body := getPullRequestBody(
		"[hugo-updater] Update Hugo to version 0.150.0",
		"https://github.com/gohugoio/hugo/releases/tag/v0.150.0",
		"Thanks to @bep for the fix",
	)
	if body != expected {
		t.Errorf("Expected %q and real %q)", expected, body)
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
