---
title: Running Classification
sidebar_position: 5
---

# Running Classification

Once your configuration file is ready, you can run FolderFlow using the
`classify` command.

---

## Basic usage

```bash
folderflow classify --config config.yaml
```
This command:
- Loads the configuration file
- Validates filters and strategies
- Scans source directories
- Classifies files into destinations
- Applies regrouping (if configured)

## Configuration flag

The configuration file must be provided explicitly.
```bash
-c, --config <path>
```

Example:
```bash
folderflow classify -c ./config.yaml
```

## Dry run mode

FolderFlow supports a dry-run mode that allows you to preview actions
without modifying the filesystem.

```bash
folderflow classify --config config.yaml --dry-run
```

In dry-run mode:
- No files are moved
- No links are created
- All decisions are logged

This is **strongly recommended** when testing new configurations.

## Output and logs

During execution, FolderFlow logs:
- Which files are matched
- Which destinations are selected
- Any errors encountered

Errors during configuration loading or classification abort execution.

## Exit behavior
- Exit code `0`: classification completed successfully
- Non-zero exit code: configuration or runtime error