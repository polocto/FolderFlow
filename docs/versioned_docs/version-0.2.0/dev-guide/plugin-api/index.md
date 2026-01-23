---
title: Plugin API
sidebar_position: 4
---

# Plugins

FolderFlow is designed around a **plugin-based architecture**.

Plugins allow you to extend or customize FolderFlow’s behavior without modifying the core engine. They encapsulate well-defined logic and participate in the file processing pipeline in a safe and deterministic way.

This section explains how plugins work, what types exist, and how to start writing your own.

---

## What is a plugin?

A plugin is a Go component that:

- Implements a public FolderFlow interface
- Is registered under a unique selector name
- Receives configuration from the YAML file
- Is loaded once at startup and reused throughout execution

Plugins are **not dynamically loaded** at runtime. They are compiled into the binary and explicitly registered.

---

## Plugin types

FolderFlow currently supports two plugin types.

### Filters

Filters decide **whether a file should be routed to a destination**.

They:
- Inspect file paths from source, metadata and data (read-only)
- Return a match decision
- Do not modify files or directories

Typical use cases:
- Filter by extension
- Filter by filename or pattern
- Filter by size or attributes
- Filter by content

📘 See: [**API Reference → Filter Interface**](./api-reference.md#filter-interface)

---

### Strategies

Strategies decide **how destination paths are computed**.

They:
- Compute a directory path for each file
- Preserve or reshape directory hierarchies and file's name
- Never interact with the filesystem

Typical use cases:
- Preserve source directory structure
- Group files by date
- Flatten directory trees

📘 See: [**API Reference → Strategy Interface**](./api-reference.md#strategy-interface)

---

## Plugin lifecycle

Plugins follow a simple and strict lifecycle:

1. Registered during program initialization
2. Instantiated when referenced in configuration
3. Configured and validated at startup
4. Used repeatedly during file processing

After configuration:
- Plugins must be immutable
- Plugins must be safe for concurrent use
- Errors during configuration abort startup

---

## Configuration model

Plugins are selected in the YAML configuration by name.

```yaml
filters:
  - name: extensions
    config:
      extensions: [".jpg", ".png"]
```
- The `type` field selects the plugin
- The `config` block is interpreted by the plugin itself
- Invalid configuration results in startup failure

Each plugin defines its own configuration schema.

## Design guarantees

All plugins must respect the following guarantees:
- No filesystem mutation
- Deterministic behavior
- No global state changes
- Safe path handling
- Clear and descriptive errors

These rules are enforced by convention and testing.

📘 See: [Best Practices](./best-practices.md)

## Getting started
If you want to write your own plugin:
1. Read the API Reference
2. Follow the Plugin Template
3. Apply the Best Practices
4. Write unit tests for edge cases