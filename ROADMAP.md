# Roadmap

This document tracks planned work for Semver-Manager, organized by command/area. Each milestone is a small, shippable increment.

---

## fetch

### Additional platforms

- [ ] GitLab
- [ ] Local git repository
- [ ] OCI registry
- [ ] ghcr.io
- [ ] npm registry
- [ ] Plain text file

### Usability

- [ ] Add progress logging on long-running fetches

### Bugs
- [] Flags `-V` and `--version` are not working

---

## filter

### New filters

- [ ] Range filter (e.g. `>=1.0.0 <2.0.0`)
- [ ] Expose the `Release` filter flag (already in `FilterArgs`)

---

## increment

### Input sources

- [ ] Accept piped input from `fetch` command
- [ ] Automated git context for build metadata

### Target Streams

- [ ] Support target stream aliases e.g. release for `*.*.*`, unmapped names for pre-release streams e.g. `alpha`, `beta`, `rc` to `*.*.*-alpha`, `*.*.*-beta`, `*.*.*-rc`

### Bugs
- [ ] Fix the following bugs
```bash
smgr increment -l patch --source-versions "0.0.1-alpha 0.0.2-alpha.0" --target-stream "*.*.*-alpha"
0.0.2-alpha.0 # should return 0.0.2-alpha.1
smgr increment -l patch --source-versions "0.0.1-alpha" --target-stream "*.*.*-alpha"
0.0.2-alpha.0 # should return 0.0.1-alpha.0
smgr increment -l patch --source-versions "" --target-stream "*.*.*-alpha"
0.0.1-alpha # should return 0.0.0-alpha or 0.0.0-alpha.0
```
---

## push

### Core implementation

- [ ] Create a tag on a target destination (GitHub, GitLab, etc.)

---

## validate

### Core implementation

- [ ] Validate a string against the Semantic Versioning 2.0.0 specification

---

## print

### Core implementation

- [ ] Create a version object and output as JSON or YAML

---

## Namespacing

- [ ] Support namespaced versions across fetch, increment, and filter

---

## Configuration

- [ ] Fix `ccs.yaml` config file loading (currently broken)
- [ ] Document env var support (`CCS_` prefix) once config is reliable
- [ ] Document flag → config → env-var precedence

---

## Distribution

- [ ] Fix `go.mod` module path to enable `go install` from GitHub
- [ ] Publish multi-platform binaries (macOS, Windows)

---

## CI / Integration

- [ ] GitHub Action for easy pipeline integration
- [ ] Usage examples for common CI providers

---

## Documentation

- [ ] Expand README description as more platforms and commands ship

---

## Long-term

- [ ] Backend server with database and API for version history management
- [ ] CLI interface improvements (TUI, interactive mode)
