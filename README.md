# Semver-Manager

[![Build](https://github.com/bluepr-nt/semver-manager/actions/workflows/on-push-to-main.yaml/badge.svg)](https://github.com/bluepr-nt/semver-manager/actions/workflows/on-push-to-main.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluepr-nt/semver-manager)](https://goreportcard.com/report/github.com/bluepr-nt/semver-manager)

A CLI tool for managing [Semantic Versioning 2.0.0](https://semver.org) compliant versions — increment, filter, and fetch version tags from GitHub repositories.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Commands](#commands)
  - [increment](#increment)
  - [filter](#filter)
  - [fetch](#fetch)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

```bash
# Increment a patch version from nothing (defaults to 0.0.1)
smgr increment --level patch
# → 0.0.1

# Increment major from existing versions
smgr increment --level major --source-versions "0.0.0,1.0.0,0.1.0"
# → 2.0.0

# Filter versions by stream and pick the highest
smgr filter --versions "1.0.0 2.0.0 1.1.0" --stream "1.*.*" --highest
# → 1.1.0

# Fetch tags from a GitHub repository
smgr fetch -o bluepr-nt -r semver-manager -t "$GITHUB_TOKEN"
```

## Installation

### Precompiled Binaries (Linux only)

1. Download the binary from the [Releases](https://github.com/bluepr-nt/semver-manager/releases) page.
2. Move the `smgr` binary to a directory in your `PATH`.

### Building from Source

Prerequisites: [Go](https://golang.org/dl/) 1.23+

```bash
cd src/ && go build -o smgr ./cmd/smgr/
```

Or with [Devbox](https://www.jetify.com/devbox) (no Go prerequisite):

```bash
devbox run build
```

## Commands

All commands support a `--dry-run` flag and standard logging flags (`-v` for verbosity).

### increment

Increment a version number (MAJOR.MINOR.PATCH) with optional pre-release support. Defaults to `0.0.1` if no source versions are provided.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--level` | `-l` | `patch` | Increment level: `major`, `minor`, `patch` (defaults to `patch` if `--target-stream` not specified) |
| `--target-stream` | `-t` | | Target stream pattern, e.g. `1.2.*` or `*.*.*-alpha.*` |
| `--source-versions` | `-s` | | Comma-separated source versions, e.g. `"0.0.0,1.0.0,1.1.0"` |

**Examples:**

```bash
# Default patch increment (no source → 0.0.1)
smgr increment

# Major increment from existing versions
smgr increment --level major --source-versions "0.0.0,1.0.0,0.1.0"
# → 2.0.0

# Pre-release increment targeting an alpha stream
smgr increment --level minor --source-versions "0.0.0,1.0.0,0.1.0" --target-stream "*.*.*-alpha.*"
# → 1.1.0-alpha.0
```

### filter

Filter a list of versions using stream patterns and/or select the highest match.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--versions` | `-V` | | Space-separated version list to filter |
| `--stream` | `-s` | | Stream pattern using `*` wildcards for any identifier |
| `--highest` | `-H` | `false` | Return only the highest version after filtering |

**Examples:**

```bash
# Filter by stream
smgr filter --versions "1.0.0 2.0.0 1.1.0" --stream "1.*.*"
# → 1.0.0 1.1.0

# Highest across all versions
smgr filter --versions "1.0.0 2.0.0 1.1.0" --highest
# → 2.0.0

# Combined: highest in a stream
smgr filter --versions "1.0.0 2.0.0 1.1.0" --stream "1.*.*" --highest
# → 1.1.0
```

<details>
<summary>Stream pattern reference</summary>

Stream patterns use `*` as a wildcard for any identifier. The absence of an identifier (or wildcard) means "no match", except for build metadata which matches anything when not specified.

| Pattern | Input versions | Result |
|---------|---------------|--------|
| `1.*.*` | 1.1.1, 2.1.1, 1.1.1+build01, 1.1.1-alpha | 1.1.1, 1.1.1+build01 |
| `*.*.*+AMD` | 1.1.1+AMD, 1.1.2, 1.1.1-alpha+AMD | 1.1.1+AMD |
| `1.0.0-Beta.*` | 1.0.0-Alpha.0, 1.0.0-Beta.0, 1.0.0-Beta.1 | 1.0.0-Beta.0, 1.0.0-Beta.1 |
| `1.0.0-Beta` | 0.1.0-Alpha, 0.1.0-Beta, 1.0.0-Beta | 1.0.0-Beta |
| `1.0.0-Beta.*` | 1.0.0-Beta.Alpha.0, 1.0.0-Beta, 1.0.0-Beta.Alpha | 1.0.0-Beta.Alpha.0, 1.0.0-Beta, 1.0.0-Beta.Alpha |
| `1.0.0-*.Beta.*` | 1.0.0-0.Alpha.0, 1.0.0-Beta.0, 1.0.0-Alpha.Beta.1 | 1.0.0-Beta.0, 1.0.0-Alpha.Beta.1 |

</details>

### fetch

Fetch semantic version tags from a GitHub repository. Automatically chains with all `filter` flags.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--owner` | `-o` | | Repository owner or organization |
| `--repo` | `-r` | | Repository name |
| `--token` | `-t` | | GitHub access token |
| `--platform` | `-p` | `github` | Platform to fetch from (currently: `github`) |
| `--stream` | `-s` | | *(from filter)* Stream pattern |
| `--highest` | `-H` | `false` | *(from filter)* Return only the highest version |
| `--versions` | `-V` | | *(from filter)* Additional versions to merge with fetched results |

**Examples:**

```bash
# Fetch all semver tags from a repo
smgr fetch -o bluepr-nt -r semver-manager -t "$GITHUB_TOKEN"

# Fetch and filter to highest in a stream
smgr fetch -o bluepr-nt -r semver-manager -t "$GITHUB_TOKEN" --stream "1.*.*" --highest
```

## Contributing

Contributions are welcomed! Please read the [Contributing Guidelines](CONTRIBUTING.md) to get started. By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

Semver-Manager is released under the [MIT License](LICENSE).
