package main

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	t.Parallel()
	expected := "netlify.toml"

	config, err := parseConfigFile("fixtures/config.toml")
	if err != nil {
		t.Fatalf("Error in parsing")
	}

	for _, repos := range config.TargetRepository {
		if repos.TargetFile != expected {
			t.Errorf("Expected %v and real %v)", expected, repos.TargetFile)
		}
	}
}

func TestParseNetlifyConfig(t *testing.T) {
	t.Parallel()
	expected := "0.83.1"

	config, err := parseNetlifyConfFile("fixtures/netlify.toml")
	if err != nil {
		t.Fatalf("Error in parsing")
	}
	real := config.Build.BuildEnv.HugoVersion
	if real != expected {
		t.Errorf("Expected %v and real %v)", expected, real)
	}

}

func TestParseNetlifyConfigDirect(t *testing.T) {
	t.Parallel()
	expected := "0.83.1"

	content := `
[build]
  command = "hugo --gc --minify -b $URL"

[build.environment]
  HUGO_VERSION = "0.83.1"
  HUGO_ENABLEGITINFO = "true"
	`

	config, err := parseNetlifyConf(content)
	if err != nil {
		t.Fatalf("Error in parsing")
	}
	real := config.Build.BuildEnv.HugoVersion
	if real != expected {
		t.Errorf("Expected %v and real %v)", expected, real)
	}

}

func TestReplaceVersion(t *testing.T) {
	input := `
[build]
  command = "hugo --gc --minify -b $URL"

[build.environment]
  HUGO_VERSION = "0.1.1"
  HUGO_ENABLEGITINFO = "true"`

	expected := `
[build]
  command = "hugo --gc --minify -b $URL"

[build.environment]
  HUGO_VERSION = "0.83.1"
  HUGO_ENABLEGITINFO = "true"`

	real := updateVersion("0.83.1", input)

	if real != expected {
		t.Errorf("Expected %v and real %v)", expected, real)
	}
}
