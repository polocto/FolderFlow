---
title: Setup
sidebar_position: 3
---

# Development Setup

This page explains how to set up FolderFlow for local development.

---

## Prerequisites

You will need:

- Go (version specified in `go.mod`)
- Node.js (v20 or later) for documentation
- Git

---

## Clone the repository

```bash
git clone https://github.com/polocto/FolderFlow.git
cd FolderFlow
```
## Build FolderFlow

```bash
go build ./cmd/ff
```

This produces the `ff` binary in the project directory.

## Run locally

```bash
./ff classify --config config.yaml
```
Use a minimal test configuration to validate behavior.

## Documentation setup

Documentation is built using Docusaurus.
```bash
cd docs
npm install
npm run start
```
This starts a local documentation server with hot reload.

## Environment notes
- FolderFlow does not require external services
- All behavior is deterministic and local
- No filesystem changes are performed during config loading

## Troubleshooting

If builds fail:
- Verify Go version matches `go.mod`
- Run `go clean -modcache`
- Ensure Node.js version is recent