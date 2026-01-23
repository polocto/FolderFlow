---
title: Filters
sidebar_position: 2
---

# Filters

Filters define which files are accepted by a destination directory.

A file must match at least one filter to be classified into a destination.

Filters do not move files and do not modify directories. They only decide whether a file matches.

---

## How filters are used

During classification:

1. FolderFlow scans files from source directories
2. Filters are evaluated for each destination
3. If a filter matches, the file is accepted
4. If no filters match, the file is ignored for that destination

Filters are evaluated independently for each destination.

---

## Built-in filters

FolderFlow currently provides the following built-in filters:

- Extensions filter

Each filter has its own configuration and behavior.

---

## Combining filters

You can define multiple filters for a destination.

When multiple filters are defined:
- A file is accepted if at least one filter matches
- Filters do not override each other
- Order does not matter

---

## Next steps

To learn how to configure filters, see:
- Extensions filter documentation
- Classification configuration examples
