---
title: Strategies
sidebar_position: 2
---

# Strategies

Strategies define how FolderFlow computes the destination directory path for files.

Strategies do not move files and do not create directories. They only compute paths.

---

## How strategies work

After a file is accepted by filters:

1. The configured strategy is invoked
2. The strategy computes a destination directory path
3. FolderFlow applies filesystem operations separately

Strategies must be deterministic and safe.

---

## Built-in strategies

FolderFlow provides the following built-in strategies:

- dirchain

Each strategy is documented in its own page.

---

## Choosing a strategy

Use a strategy when you want to:
- Preserve source directory hierarchy
- Control how files are organized in destinations
- Prevent filename collisions

If no strategy is configured, FolderFlow cannot determine where files should be placed.

---

## Next steps

See the dirchain strategy documentation for details and examples.
