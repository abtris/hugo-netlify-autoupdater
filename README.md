[![Build Status](https://github.com/abtris/hugo-netlify-autoupdater/actions/workflows/go.yaml/badge.svg)](https://github.com/abtris/hugo-netlify-autoupdater/actions)
# Hugo Netlify Autoupdater

- [ ] listen on event publish release `gohugoio/hugo`
  - cron in 1st iteration
- [x] compare current deployed version of Hugo in all blogs
- [x] create PR's for update version
- [ ] merge if all passed
- [ ] more settings move into config file

## Dependencies

- [go-github library](https://github.com/google/go-github)
- [toml library](https://github.com/BurntSushi/toml)
- [go-version](https://github.com/hashicorp/go-version)

## Install

```sh
git clone https://github.com/abtris/hugo-netlify-autoupdater.git
cd hugo-netlify-autoupdater
go mod download
```

## Config

Configuration is in `config.toml` file.

```toml
source_repo_releases = "gohugoio/hugo"

[[target_repos]]
  repo = "owner/repo"
  target_file = "netlify.toml"
  target_variable = "HUGO_VERSION"
  branch = "master"
```
## Run

```sh
make run
```

## License

MIT
