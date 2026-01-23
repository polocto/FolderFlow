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
type Context interface {
	PathFromSource() string
	DstDir() string
	Info() fs.FileInfo
}

type Strategy interface {
	FinalDirPath(ctx Context) (string, error)
	Selector() string
	LoadConfig(config map[string]interface{}) error
}
```

### `FinalDirPath`

```go
func FinalDirPath(ctx Context) (string, error)
```

**Purpose**

Computes the final directory where a file should be placed.

**Rules**
- MUST NOT modify the filesystem
- MUST be deterministic
- MUST return a directory path (not a filename)

**Parameters**
- `ctx` that should return all usefull information for strategy computation

### `Selector`

```go
func Selector() string
```

Returns a **unique** identifier for the strategy.

Example values:
- "date"
- "dirchain"
- "custom"

This value is used in configuration files.

### `LoadConfig`

```go
func LoadConfig(config map[string]interface{}) error
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
type Context interface {
	Info() fs.FileInfo
	WithInput(fn func(r io.Reader) error) error
	WithInputLimited(maxBytes int64, fn func(r io.Reader) error) error
	ReadChunks(chunkSize int, fn func([]byte) error) error
}

type Filter interface {
	Match(ctx Context) (bool, error)
	Selector() string
	LoadConfig(config map[string]interface{}) error
}
```

### `Match`

```go
func Match(ctx Context) (bool, error)
```

Returns true if the file matches the filter criteria.

Filters:
- `ctx` that should return all usefull information for filter computation

### `Selector`

Same semantics as strategies: a unique identifier used in configuration.

### `LoadConfig`

Loads filter-specific configuration.

## Plugin Registration

Strategies and filters must be registerd before use.

Registration happens in an `init()` function.

Force import in `main.go`

```go
import _ "github.com/yourorg/folderflow/plugins/myfilter"
```

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