package main

import (
	"context"
	"testing"

	"github.com/google/go-github/v35/github"
)

// Expecting test fail after new release (remove later)
func TestGetCurrentHugoVersion(t *testing.T) {
	expected := "0.83.1"

	var client *github.Client
	var ctx = context.Background()

	client = github.NewClient(nil)

	real, err := getCurrentHugoVersion(ctx, client)
	if err != nil {
		t.Fatalf("Get error %v", err)
	}
	if real != expected {
		t.Errorf("Expected %v and real %v)", expected, real)
	}
}

func TestGetCurrentDeployedVersion(t *testing.T) {
	expected := "0.83.1"

	var client *github.Client
	var ctx = context.Background()

	client = github.NewClient(nil)

	real, err := getCurrentDeployedVersion(ctx, client, "abtris", "www.prskavec.net", "netlify.toml")
	if err != nil {
		t.Fatalf("Get error %v", err)
	}

	if real != expected {
		t.Errorf("Expected %v and real %v)", expected, real)
	}
}

func TestIsNewVersion(t *testing.T) {
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
