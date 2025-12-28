---
title: Strategies
sidebar_position: 6
---

# Strategies

Strategies define how destination directory paths are computed for files.

They do not move files or create directories. They only compute paths.

---

## How strategies work

For each accepted file:
1. FolderFlow invokes the configured strategy
2. The strategy computes a destination directory
3. FolderFlow applies filesystem operations separately

Strategies must be deterministic and safe.

---

## Built-in strategies

FolderFlow provides the following built-in strategies:

- dirchain

Each strategy is documented in its own page.

---

## Writing custom strategies

Custom strategies:
- Implement the Strategy interface
- Receive configuration via YAML
- Must never touch the filesystem


---


## Strategy Template

```go
package mystrategy

import (
    "io/fs"
)

type MyStrategy struct {
    // configuration fields
}

func (s *MyStrategy) Selector() string {
    return "my-strategy"
}

func (s *MyStrategy) LoadConfig(config map[string]interface{}) error {
    // parse and validate config
    return nil
}

func (s *MyStrategy) FinalDirPath(
    srcDir, destDir, filePath string,
    info fs.FileInfo,
) (string, error) {
    // compute destination directory
    return destDir, nil
}

func init() {
    strategy.RegisterStrategy("my-strategy", func() strategy.Strategy {
        return &MyStrategy{}
    })
}

```