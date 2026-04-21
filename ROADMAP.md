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
