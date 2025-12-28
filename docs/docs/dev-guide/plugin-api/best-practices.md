---
title: Best Practice
sidebar_position: 3
---

# Best Practices

This guide outlines recommended practices when using or extending FolderFlow.

---

## Do Not Touch the Filesystem in Plugins

Filters and strategies **must not**:
- Move files
- Create directories
- Delete files

They should **only compute decisions**.

---

## Make Plugins Deterministic

Given the same input:
- Filters should return the same result
- Strategies should return the same path

Avoid:
- Randomness
- Time-based logic (unless explicitly configured)

---

## Validate Configuration Early

Fail fast in `LoadConfig`:
- Check required fields
- Validate types
- Provide meaningful error messages

---

## Prefer Composition Over Complexity

Small, composable plugins are easier to:
- Test
- Reuse
- Debug

Avoid “do-everything” strategies.

---

## Log Carefully

Plugins may log for debugging, but:
- Avoid excessive logging
- Never log sensitive paths by default

---

## Version Your Plugins Carefully

Changes to plugin behavior may affect existing workflows.

When breaking behavior:
- Bump plugin version
- Document changes clearly