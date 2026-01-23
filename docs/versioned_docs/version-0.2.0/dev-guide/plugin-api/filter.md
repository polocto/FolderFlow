---
title: Filters
sidebar_position: 5
---

# Filters

Filters decide whether a file should be routed to a destination.

They are evaluated during classification and return a simple yes or no decision.
Filters never modify the filesystem.

---

## How filters work

For each file:
1. FolderFlow invokes all configured filters
2. If at least one filter matches, the file is accepted
3. If no filter matches, the file is ignored for that destination

Filters are evaluated independently and must be deterministic.

---

## Built-in filters

FolderFlow provides the following built-in filters:

- Extensions filter
- Regex filter

Each filter is documented in its own page.

---

## Writing custom filters

If built-in filters are not sufficient, you can write your own.

Custom filters:
- Implement the Filter interface
- Are configured via YAML
- Are registered at startup


---

## Filter Template

```go
package myfilter

import (
    "io/fs"
    "github.com/polocto/FolderFlow/pkg/ffplugin/filter"
)

type MyFilter struct {
    // configuration fields
}

func (f *MyFilter) Selector() string {
    return "my-filter"
}

func (f *MyFilter) LoadConfig(config map[string]interface{}) error {
    return nil
}

func (f *MyFilter) Match(path string, info fs.FileInfo) (bool, error) {
    return true, nil
}

func init(){
    filter.RegisterFilter("my-filter", func() filter.Filter {
        return &MyFilter{}
    })
}
```

## Import package for build

Force import in `main.go`

```go
import _ "github.com/yourorg/folderflow/plugins/myfilter"
```