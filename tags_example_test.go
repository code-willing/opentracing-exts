package trace_test

import (
	"net"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/code-willing/trace"
)

func ExampleRPCTags_client() {
	// Start a new span with RPC client tags set.
	span := opentracing.StartSpan("name", trace.RPCTags{
		Kind:        ext.SpanKindRPCClientEnum,
		PeerAddr:    "http://service.name.io/",
		PeerService: "service",
	})
	defer span.Finish()
}

func ExampleRPCTags_server() {
	// Start a new span with RPC server tags set.
	span := opentracing.StartSpan("name", trace.RPCTags{
		Kind:        ext.SpanKindRPCServerEnum,
		PeerAddr:    "http://service.name.io/",
		PeerService: "service",
		PeerIPv4:    net.IPv4(127, 0, 0, 1),
		PeerPort:    8080,
	})
	defer span.Finish()
}

func ExampleDBTags_sql() {
	// Start a new span with database client tags set.
	span := opentracing.StartSpan("name", trace.DBTags{
		Type:      "sql",
		Instance:  "test",
		User:      "username",
		Statement: "SELECT * FROM test",
	})
	defer span.Finish()
}

func ExampleDBTags_nosql() {
	// Start a new span with database client tags set.
	span := opentracing.StartSpan("name", trace.DBTags{
		Type:      "redis",
		Instance:  "test",
		User:      "username",
		Statement: "SET mykey 'WuValue'",
	})
	defer span.Finish()
}

func ExampleHTTPTags() {
	// Start a new span with HTTP tags set.
	span := opentracing.StartSpan("name", trace.HTTPTags{
		Method:     http.MethodGet,
		URL:        "http://example.com/test?foo=bar",
		StatusCode: http.StatusOK,
	})
	defer span.Finish()
}
