---
title: Configuration
sidebar_position: 1
---
# Configuration File Format

FolderFlow uses a **YAML configuration file** to define its entire behavior.

This file controls:
- Source directories to scan
- Destination directories
- File filters
- Directory structure strategies
- Regrouping behavior
- Parallelism


The configuration file is written in YAML and typically named:

```text
config.yaml
```
Once your configuration is ready, see [**Running Classification**](./run.md) to execute FolderFlow.

## High-Level Structure

```yaml
source_dirs: []
dest_dirs: []
regroup: {}
max_workers: 0
```
Each section is described in detail below.

## Source Directories

```yaml
source_dirs:
  - "./source/path"
```
* One or more directories to scan recursively
* All files inside are evaluated against destination filters
* Directory hierarchy is preserved during processing

## Destination Directories

```yaml
dest_dirs:
  - name: "images"
    path: "./destination/images"
    filters: []
    strategy: {}
```

Each destination defines:
* What files it accepts
* Where those files are stored
* How directory structure is recreated


## Filters

Filters control which files belong to a destination.

### `extensions`

```yaml
filters:
  - name: "extensions"
    config:
      extensions: [".jpg", ".png"]
```
* Matches files by extension
* Multiple extensions allowed
* Case-insensitive matching

A file must match at least one filter to be routed to the destination.

## Strategy
Strategies define how directory structure is rebuilt in the destination.

### `dirchain`

```yaml
strategy:
  name: "dirchain"
```
* Recreates the full relative directory path
* Preserves source hierarchy
* Prevents filename collisions

## Regroup

```yaml
regroup:
  path: "./regrouped"
  mode: hardlink
```
The regroup section creates a central access directory containing links to all processed files.
* Files are not duplicated
* Useful for indexing, previewing, or downstream processing

## Parallelism
```yaml
max_workers: 0
```
* `0` → automatically determined based on CPU cores
* `> 0` → fixed number of concurrent workers

Higher values increase speed but also disk and CPU usage.

