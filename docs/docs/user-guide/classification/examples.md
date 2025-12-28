---
title: Examples
sidebar_position: 3
---

# Example 1 — Organizing Mixed Files

```yaml
source_dirs:
  - "./downloads"

dest_dirs:
  - name: "documents"
    path: "./sorted/documents"
    filters:
      - name: "extensions"
        config:
          extensions: [".pdf", ".txt", ".md"]
    strategy:
      name: "dirchain"

  - name: "images"
    path: "./sorted/images"
    filters:
      - name: "extensions"
        config:
          extensions: [".jpg", ".png"]
    strategy:
      name: "dirchain"

max_workers: 0
```

# Example 2 - Media Collection With Regroup

```yaml
source_dirs:
  - "./media"

dest_dirs:
  - name: "videos"
    path: "./library/videos"
    filters:
      - name: "extensions"
        config:
          extensions: [".mp4", ".mov"]
    strategy:
      name: "dirchain"

regroup:
  path: "./library/all"
  mode: hardlink
```

# Example 3 - Large Dataset, Limited Parallelism

```yaml
source_dirs:
  - "./dataset"

dest_dirs:
  - name: "archives"
    path: "./organized/archives"
    filters:
      - name: "extensions"
        config:
          extensions: [".zip", ".tar.gz"]
    strategy:
      name: "dirchain"

max_workers: 4
```