# efmt

[![tag](https://img.shields.io/github/tag/vasilesk/efmt.svg)](https://github.com/vasilesk/efmt/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-%23007d9c)
[![Go Reference](https://pkg.go.dev/badge/github.com/vasilesk/efmt.svg)](https://pkg.go.dev/github.com/vasilesk/efmt)
![Build Status](https://github.com/vasilesk/efmt/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/vasilesk/efmt)](https://goreportcard.com/report/github.com/vasilesk/efmt)
[![Coverage](https://img.shields.io/codecov/c/github/vasilesk/efmt)](https://codecov.io/gh/vasilesk/efmt)
[![License](https://img.shields.io/github/license/vasilesk/efmt)](./LICENSE)

`efmt` is a Go library for enhanced error handling that allows you to attach structured metadata (key-value pairs) to errors. This helps to enrich error context while maintaining compatibility with Go's standard error interface.

## Installation

```bash
go get github.com/vasilesk/efmt@latest
```

## Features

- Create errors with structured metadata
- Wrap existing errors with additional context
- Extract metadata from errors as key-value pairs
- Collect metadata from the whole error tree, including `fmt.Errorf("%w")` and `errors.Join`
- Pass metadata to `log/slog`
- Fully compatible with Go's standard error interface and unwrap mechanism
- No dependencies outside the standard library

## Usage Examples

### Creating a new error with metadata

```go
import "github.com/vasilesk/efmt"

func validateUser(user User) error {
    if user.Age < 18 {
        return efmt.New("user is underage",
            efmt.KV("user_id", user.ID),
            efmt.KV("age", user.Age),
            efmt.KV("required_age", 18),
        )
    }
    return nil
}
```

### Wrapping an existing error

```go
import "github.com/vasilesk/efmt"

func processRequest(req Request) error {
    user, err := fetchUser(req.UserID)
    if err != nil {
        return efmt.Wrap(err, "failed to fetch user",
            efmt.KV("request_id", req.ID),
            efmt.KV("user_id", req.UserID),
        )
    }
    // Process user...
    return nil
}
```

### Extracting metadata from errors

```go
import (
    "fmt"
    "github.com/vasilesk/efmt"
)

func handleError(err error) {
    // Get all key-value pairs
    allValues := efmt.Values(err)
    fmt.Printf("Error context: %+v\n", allValues)

    // Get a specific value
    if requestID, ok := efmt.Value[string](err, "request_id"); ok {
        fmt.Printf("Request ID: %s\n", requestID)
    }

    // Get all key-value pairs as a slice
    pairs := efmt.ValuePairs(err)
    for _, kv := range pairs {
        fmt.Printf("%s: %v\n", kv.Key, kv.Value)
    }
}
```

### Working with error chains

```go
import "github.com/vasilesk/efmt"

func processData() error {
    data, err := fetchData()
    if err != nil {
        return efmt.Wrap(err, "data fetch failed",
            efmt.KV("timestamp", time.Now()),
        )
    }

    err = validateData(data)
    if err != nil {
        return efmt.Wrap(err, "data validation failed",
            efmt.KV("data_size", len(data)),
        )
    }

    return nil
}

// Later when handling the error:
func handleProcessingError(err error) {
    // This will extract all metadata from the entire error chain
    allValues := efmt.Values(err)
    fmt.Printf("Complete error context: %+v\n", allValues)
}
```

### Logging with slog

`efmt.Attrs` returns the metadata as `slog` attributes in a stable order:

```go
import (
    "context"
    "log/slog"

    "github.com/vasilesk/efmt"
)

func logError(ctx context.Context, logger *slog.Logger, err error) {
    attrs := append(efmt.Attrs(err), slog.Any("error", err))
    logger.LogAttrs(ctx, slog.LevelError, "request failed", attrs...)
}
```

## Key precedence

When the same key is set at several levels of an error tree, the outermost value wins:

```go
inner := efmt.New("parse failed", efmt.KV("stage", "parse"), efmt.KV("line", 7))
outer := efmt.Wrap(inner, "load failed", efmt.KV("stage", "load"))

efmt.Values(outer) // map[line:7 stage:load]
```

`ValuePairs` returns all pairs, duplicates included, in traversal order: an error comes
before the errors it wraps, and errors combined with `errors.Join` are visited in order.
`Values`, `Value` and `Attrs` all use the first occurrence of a key in that order.
`Value[T]` returns `false` if that occurrence is not of type `T`; it does not fall back
to a shadowed inner value.

## Compatibility with the errors package

`efmt` errors work with `errors.Is`, `errors.As` and `errors.Unwrap`, and metadata is
collected through errors wrapped by other means:

```go
import (
    "errors"
    "fmt"

    "github.com/vasilesk/efmt"
)

err := fmt.Errorf("handler: %w", efmt.Wrap(ErrNotFound, "lookup failed", efmt.KV("id", 42)))

errors.Is(err, ErrNotFound) // true
efmt.Values(err)            // map[id:42]
```

## License

See the [LICENSE](LICENSE) file for details.
