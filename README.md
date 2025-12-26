# FolderFlow

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/github/license/polocto/FolderFlow)](LICENSE)
[![CI](https://github.com/polocto/FolderFlow/actions/workflows/go.yaml/badge.svg)](https://github.com/polocto/FolderFlow/actions/workflows/go.yaml)
[![Last Commit](https://img.shields.io/github/last-commit/polocto/FolderFlow/dev)](https://github.com/polocto/FolderFlow/commits/dev)
[![Last Release](https://img.shields.io/github/v/release/polocto/FolderFlow)](https://github.com/polocto/FolderFlow/releases)

FolderFlow is a **Go-based command-line application** for classifying and organizing files from one or more source directories into structured destination directories, based on a **YAML configuration file**.

It scans directories recursively, applies filters to files, preserves directory structure, and moves files accordingly.  
An optional *regroup* directory can be created to centralize all processed files using filesystem links.

---

## Key Concepts

FolderFlow is **configuration-driven**.

The CLI is intentionally minimal:  
most behavior is defined in a YAML file, not via command-line arguments.

At a high level, FolderFlow:
1. Walks through source directories
2. Matches files against destination filters
3. Moves matching files into destination directories
4. Preserves directory hierarchy
5. Optionally regroups processed files using links
6. Processes files concurrently

---

## Installation

### Build from source

```bash
git clone https://github.com/polocto/FolderFlow.git
cd FolderFlow
go build
```

This produces a `FolderFlow` (or `FolderFlow.exe`) binary.

## Usage

FolderFlow is executed through the `classify` command.

```bash
folderflow classify --file config.yaml
```

### Globals flags

Available on all commands:
- `--verbose`, `-v`
    
    Enable verbose logging
- `--debug`
    
    Enable debug-level logs
- `--dry-run`
    
    Perform a dry run (no filesystem changes)

## Configuration file

```yaml
source_dirs:
  - "./testdata/source"

dest_dirs:
  - name: "documents"
    path: "./testdata/destination/documents"
    filters:
      - name: "extensions"
        config:
          extensions: [".txt", ".pdf", ".md"]
    strategy:
      name: "dirchain"

  - name: "images"
    path: "./testdata/destination/images"
    filters:
      - name: "extensions"
        config:
          extensions: [".jpg", ".png", ".gif"]
    strategy:
      name: "dirchain"

regroup:
  path: "./testdata/regrouped"
  mode: hardlink

max_workers: 0
```

### Configuration Reference

#### `source_dirs`
```yaml
source_dirs:
  - "/path/to/source"
```
- List of directories to scan recursively
- Non-existing directories are skipped
- Certain directories are always skipped internally:
    - .git
    - node_modules
- A source directory cannot be the same as the regroup path

#### `dest_dirs`
Defines where files are moved when they match filters.
```yaml
dest_dirs:
  - name: "images"
    path: "/destination/images"
    filters: [...]
    strategy: ...
```
**Fields**
|Field|Description|
|:-|:-|
|`name`|Logical name (used for clarity and logging)|
|`path`|Destination root directory|
|`filters`|List of filters that must **all match** for the file to be routed here|
|`strategy`|Controls how the destination path is built|

#### Filters

Currently implemented filter:

`extensions`
```yaml
filters:
  - name: "extensions"
    config:
      extensions: [".jpg", ".png"]
```
- Matches files based on file extension
- Comparison is done per file
- All filters must match for a destination to apply

#### Strategy

Currently implemented strategy:

`dirchain`
```yaml
strategy:
  name: "dirchain"
```
- Recreates the directory structure from the source directory
- Prevents filename flattening
- Helps avoid collisions and keeps context

#### Regroup
Optional section to centralize all processed files.
```yaml
regroup:
  path: "./regrouped"
  mode: hardlink
```
- A link is created after the file is moved
- The regroup directory contains one link per processed file
- No file duplication

Supported modes:

|Mode|Description|
|:-|:-|
|`symlink`||
|`hardlink`||
|`copy`||

If the regroup path matches a source directory, it is skipped to avoid loops.

#### `max_workers`
```yaml
max_workers: 0
```
- Controls concurrency
- `0` means automatic worker count
- Affects file processing, not directory walking

## Safety Features

- Skips .git and node_modules directories
- Prevents source/destination overlap
- Supports dry-run mode
- Recovers from worker panics
- Logs all errors without stopping the entire run

## Project Status

FolderFlow is **under active development**.

### Development Status
- Core classification pipeline is implemented and functional
- Configuration-driven workflow is stable
- Concurrency, error handling, and statistics collection are in place
- CLI structure is established using Cobra

Some command descriptions and help texts are still placeholders and will be refined as the project matures.

### CI / Pipeline
- Continuous Integration is enabled on the repository
- Builds and tests are executed automatically on each push and pull request
- Pipeline status reflects the current stability of the codebase

Refer to the repository’s **Actions** tab for the latest pipeline runs and results.

### Repository Activity
- The project is actively maintained
- New features, refactors, and documentation updates are ongoing
- Breaking changes may occur until a stable release is tagged

Users are encouraged to:
- Watch the repository for updates
- Report issues or unexpected behavior
- Contribute via pull requests or discussions

## Versioning & Changelog

FolderFlow uses **explicit versioning** to track changes over time.

All notable changes are documented in the repository changelog:

- 📄 [`CHANGELOG.md`](CHANGELOG.md)

Until a stable `v1.0.0` release:
- Versions may introduce breaking changes
- Configuration format may evolve
- Users should consult the changelog before upgrading

## License

This project is licensed under the MIT License.

See the license file for details:

- 📄 [`LICENSE`](LICENSE)
