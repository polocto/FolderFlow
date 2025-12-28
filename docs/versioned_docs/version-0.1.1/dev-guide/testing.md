---
title: Testing
sidebar_position: 4
---

# Testing

Testing is critical to maintaining FolderFlow’s safety guarantees.

---

## Running tests

Run all tests:

```bash
go test ./...
```
Run with race detection:
```bash
go test -race ./...
```

## What to test

Contributors should test:
- Filters independently
- Strategies independently
- Path edge cases
- Invalid configuration handling

## Filter tests
Filters should be tested for:
- Positive matches
- Negative matches
- Case-insensitivity
- Invalid configuration

Filters must never modify the filesystem.

## Strategy tests

Strategies must be tested for:
- Correct path computation
- Root-level files
- Nested directories
- Path traversal prevention
- Directory input rejection


## CI expectations

All tests must pass in CI:
- Unit tests
- Lint checks
- Docs build

PRs with failing tests will not be merged.

## Writing new tests

Prefer:
- Table-driven tests
- Explicit expected errors
- Clear test names

Avoid:
- Integration-heavy tests
- Global state