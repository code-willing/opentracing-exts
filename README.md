# opentracing-exts
[![GoDoc](https://godoc.org/github.com/code-willing/opentracing-exts?status.svg)](https://godoc.org/github.com/code-willing/opentracing-exts)
[![Go Report Card](https://goreportcard.com/badge/github.com/code-willing/opentracing-exts)](https://goreportcard.com/report/github.com/code-willing/opentracing-exts)

This package provides opentracing span options and span logging utility types
and functions.

```
go get -u github.com/code-willing/opentracing-exts
```

```go
import (
    otexts "github.com/code-willing/opentracing-exts"
)
```

## Logging

Easily log span errors and set the correct span tags and log fields.

```go
func example() (err error) {
    span := opentracing.StartSpan("name")
    defer func() {
        if err != nil {
            otexts.LogError(span, err)
        }
        span.Finish()
    }()
}
```

Log an error with extra log fields.

```go
func example(a, b int) (err error) {
    span := opentracing.StartSpan("name")
    defer func() {
        if err != nil {
            otexts.LogErrorWithFields(span, err, map[string]interface{}{
                "a": a,
                "b": b,
            })
        }
        span.Finish()
    }()
}
```

JSON marshal the log field values.

```go

type jsonThing struct {
    a int `json:"a"`
    b int `json:"a"`
}

func example(t jsonThing) (err error) {
    span := opentracing.StartSpan("name")
    defer func() {
        if err != nil {
            otexts.LogErrorWithFields(span, err, otexts.LogFields{
                "my-thing": t,
            }.Encode())
        }
        span.Finish()
    }()
}
```

## Span Options

Start spans with specific client/server tags set:

```go
// Start a span with RPC client tags set.
span := opentracing.StartSpan("name", otexts.RPCTags{
    Kind:        ext.SpanKindRPCClientEnum,
    PeerAddr:    "http://service.name.io/",
    PeerService: "service",
})

// Start a span with RPC server tags set.
span := opentracing.StartSpan("name", otexts.RPCTags{
    Kind:        ext.SpanKindRPCClientEnum,
    PeerAddr:    "http://service.name.io/",
    PeerService: "service",
})

// Start a new span with SQL database client tags set.
span := opentracing.StartSpan("name", trace.DBTags{
    Type:      "sql",
    Instance:  "test",
    User:      "username",
    Statement: "SELECT * FROM test",
})

// Start a new span with NoSQL database client tags set.
span := opentracing.StartSpan("name", trace.DBTags{
	Type:      "redis",
	Instance:  "test",
	User:      "username",
	Statement: "SET mykey 'WuValue'",
})

// Start a new span with HTTP tags set.
span := opentracing.StartSpan("name", trace.HTTPTags{
	Method:     http.MethodGet,
	URL:        "http://example.com/test?foo=bar",
	StatusCode: http.StatusOK,
})
```

Set span tags for events:

```go
func execDB(stmt database.Stmt) (err error) {
    span := opentracing.StartSpan("name")
    defer func() {
        if err != nil {
            otexts.SetDBTags(span, otexts.DBTags{
                Type:      "sql",
                Instance:  "test",
                User:      "username",
                Statement: stmt.String(),
            })
            otexts.LogErrorWithFields(span, err, map[string]interface{}{
                "a": a,
                "b": b,
            })
        }
        span.Finish()
    }()
}
```
