---
title: dirchain
sidebar_position: 1
---

# Dirchain Strategy

The dirchain strategy preserves the relative directory structure of source files when placing them into destination directories.

It is the simplest and safest strategy for most use cases.

---

## Selector name

dirchain

---
## Configuration
The dirchain strategy does not require any configuration.

You only need to specify its name.
```yaml
strategy:
  name: "dirchain"
```
---
## Behavior

For each file:
- The directory structure relative to the source directory is preserved
- The file is placed into the corresponding destination subdirectory

Files located at the root of a source directory are placed at the root of the destination directory.

### Example

Source structure:
- source/docs/report.pdf
- source/images/photo.jpg

Destination structure using dirchain:
- destination/docs/report.pdf
- destination/images/photo.jpg

---

## Notes

The dirchain strategy guarantees that:
- Directory traversal is prevented
- Computed paths always remain inside the destination directory
- Source root files map to destination root
- No filesystem operations are performed by the strategy itself