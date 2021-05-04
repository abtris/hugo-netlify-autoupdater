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
