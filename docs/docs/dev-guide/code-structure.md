---
title: Code Structure
sidebar_position: 2
---

# Code Structure

This page explains the high-level structure of the FolderFlow codebase.

---

## Top-level layout

```text
cmd/
  ff/            # CLI entry point
filter/          # Filter plugins and registry
strategy/        # Strategy plugins and registry
internal/        # Core engine logic
docs/            # Documentation (Docusaurus)
```

### `cmd/ff`
Contains:
- CLI parsing
- Configuration loading
- Application bootstrap

This layer should remain thin.

### `filter`
Contains:
- Filter interface
- Built-in filters
- Filter registry

Filters are:
- Pure
- Stateless after configuration
- Non-destructive

### `strategy`

Contains:
- Strategy interface
- Built-in strategies
- Strategy registry

Strategies:
Compute destination paths
- Never touch the filesystem
- Enforce path safety

## Registries

Filters and strategies are registered using explicit registration functions.

This design:
- Avoids reflection
- Enables compile-time safety
- Makes plugins explicit and testable

## Design philosophy
- Clear separation of concerns
- Deterministic behavior
- Defensive path handling
- Explicit errors over silent failure