# Semver-Manager

Semver-Manager is a command-line interface (CLI) tool that streamlines semantic versioning management for developers and seamlessly integrates with popular Git repositories or registry platforms. With Semver-Manager, you can rapidly generate new versions, maintain your version history, and ensure compliance with the Semantic Versioning 2.0.0 specification, all while effortlessly working alongside your preferred platform.
## Features

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Features
### Implemented
- Increment  
  The increment command allows to increment a version's MAJOR, MINOR or PATCH number  

- Fetch  
  The fetch command allows to fetch semantic versions from multiple platforms and filter them 
  The fetch command is automatically chained with the filter command  
  Platforms:
  - Github

- Filter
  - Highest: returns the highest semver found, is always run after all other filters
  - Stream: returns all version matching the requested identifiers and accepts any identifiers when a '*' wildcard is specified. The absence of an identifier or wildcard equals a no match, excepts for build metadata which is always a match when not specified in stream filter pattern.
    Examples:
    - Pattern: 1.\*.*
      Versions: 1.1.1, 2.1.1, 1.1.1+build01, 1.1.1-alpha
      Result: 1.1.1, 1.1.1+build01
    - Pattern: \*.\*.*+AMD
      Versions: 1.1.1+AMD, 1.1.2, 1.1.1-alpha+AMD
      Result: 1.1.1+AMD
### To do
#### Fetch
- Implement additional platforms: gitlab, local git repository, oci repository, ghrc.io, npm, text file

#### Filter
- range

#### increment
- increment version MAJOR, MINOR, PATCH, defaults to latest number
- Increment major, minor, or patch versions with a single command

#### Push
- create tag <tag> on `<destination> `

#### Print
- create version object and return json or yaml

#### Validate
- Validate semantic versioning

#### CI interface
- Easy integration with CI/CD pipelines (github action)

#### CMD interface

#### Backend server with database and api
- Manage your project's version history

## Installation

To install Semver-Manager, you can download the binary for your platform from the [Releases](https://github.com/13013SwagR/semver-manager/releases) page or build it from source.

### Precompiled Binaries (Linux only)

1. Download the appropriate binary for your platform from the [Releases](https://github.com/13013SwagR/semver-manager/releases) page.
2. Extract the archive and move the `semver-manager` binary to a directory in your system's `PATH`.

### Building from Source

Prerequisites:

- [Go](https://golang.org/dl/) 1.20+ installed and configured

```sh
go install github.com/13013SwagR/semver-manager/smgr@<git_ref>
```

## Usage

```
Manage Semantic Versioning compliant versions and integrate with popular or registry platform to facilitate the task.

Usage:
  smgr [flags]
  smgr [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  fetch       Fetch semver tags from a registry or repository.
  help        Help about any command

Flags:
      --add_dir_header                   If true, adds the file directory to the header of the log messages
      --alsologtostderr                  log to standard error as well as files (no effect when -logtostderr=true)
      --dry-run                          Execute the command in dry-run mode
  -h, --help                             help for smgr
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory (no effect when -logtostderr=true)
      --log_file string                  If non-empty, use this log file (no effect when -logtostderr=true)
      --log_file_max_size uint           Defines the maximum size a log file can grow to (no effect when -logtostderr=true). Unit is megabytes. If the value is 0, the maximum file size is unlimited. (default 1800)
      --logtostderr                      log to standard error instead of files
      --one_output                       If true, only write logs to their native severity level (vs also writing to each lower severity level; no effect when -logtostderr=true)
      --skip_headers                     If true, avoid header prefixes in the log messages
      --skip_log_headers                 If true, avoid headers when opening log files (no effect when -logtostderr=true)
      --stderrthreshold severity         logs at or above this threshold go to stderr when writing to files and stderr (no effect when -logtostderr=true or -alsologtostderr=false) (default 2)
  -v, --v Level                          number for the log level verbosity
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

Use "smgr [command] --help" for more information about a command.
```

## Examples

Initialize a new project with a default version file:

```sh
smgr fetch -r semver-manager -o 13013SwagR -t <github_token>
```

## Configuration
Semver-Manager looks for a file in the current directory named `ccs.yaml`

All flags are available as configuration entries, for example:
```
TOKEN: <github_token>
REPO: semver-manager
OWNER: SMARTeacher
```

All flags are also available environment variables with the `CCS_` prefix, for example:   
`TOKEN=<github_token>` `REPO=semver-manager` `OWNER=SMARTeacher`
## Contributing

Contributions to Semver-Manager are welcomed and appreciated! Please read the [Contributing Guidelines](CONTRIBUTING.md) to get started. By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

Semver-Manager is released under the [MIT License](LICENSE).
