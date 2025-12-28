---
title: API Reference
sidebar_position: 2
---


# API Reference

This section documents the public Go interfaces used to extend FolderFlow.

---

## Strategy Interface

Strategies define **how destination paths are computed** for files.

```go
type Strategy interface {
    FinalDirPath(srcDir, destDir, filePath string, info fs.FileInfo) (string, error)
    Selector() string
    LoadConfig(config map[string]interface{}) error
}
```

### `FinalDirPath`

```go
FinalDirPath(
  srcDir string,
  destDir string,
  filePath string,
  info fs.FileInfo,
) (string, error)
```

**Purpose**

Computes the final directory where a file should be placed.

**Rules**
- MUST NOT modify the filesystem
- MUST be deterministic
- MUST return a directory path (not a filename)

**Parameters**
- `srcDir`: root source directory
- `destDir`: root destination directory
- `filePath`: path to the file being processed
- `info`: file metadata

### `Selector`

```go
Selector() string
```

Returns a **unique** identifier for the strategy.

Example values:
- "date"
- "dirchain"
- "custom"

This value is used in configuration files.

### `LoadConfig`

```go
LoadConfig(config map[string]interface{}) error
```

LoadConfig receives the content of the `config` block from the YAML configuration verbatim.
Plugins are responsible for validating and applying their configuration.


Implementations should:
- Validate required fields
- Apply defaults when possible
- Return descriptive errors

## Filter Interface

Filters define which files should be processed.

```go
type Filter interface {
    Match(path string, info fs.FileInfo) (bool, error)
    Selector() string
    LoadConfig(config map[string]interface{}) error
}
```

### `Match`

```go
Match(path string, info fs.FileInfo) (bool, error)
```

Returns true if the file matches the filter criteria.

Filters:
- Should be fast
- Should not modify the filesystem
- Can inspect file and its metadata

### `Selector`

Same semantics as strategies: a unique identifier used in configuration.

### `LoadConfig`

Loads filter-specific configuration.

## Plugin Registration

Strategies and filters must be registerd before use.

Registration happens in an `init()` function.

### Registering a Strategy

```go
func init(){
    strategy.RegisterStrategy("my-strategy", func() strategy.Strategy {
        return &MyStrategy{}
    })
}
```

### Registering a Filter

```go
func init(){
    filter.RegisterFilter("my-filter", func() filter.Filter {
        return &MyFilter{}
    })
}
```