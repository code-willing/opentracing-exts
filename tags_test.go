package trace_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"

	"github.com/code-willing/trace"
)

func init() {
	opentracing.SetGlobalTracer(mocktracer.New())
}

func TestRPCTags_Apply(t *testing.T) {
	tt := []struct {
		name string
		tags trace.RPCTags
	}{
		{
			name: "all tags",
			tags: trace.RPCTags{
				Kind:         ext.SpanKindRPCClientEnum,
				PeerAddr:     "http://internal.service.io/",
				PeerHostname: "internal.service.io",
				PeerIPv4:     net.IPv4(127, 0, 0, 1),
				PeerIPv6:     net.IPv6zero,
				PeerPort:     8080,
				PeerService:  "service",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test", tc.tags).(*mocktracer.MockSpan)
			span.Finish()
			ensureRPCTagsSet(t, tc.tags, span.Tags())
		})
	}
}

func TestSetRPCTags(t *testing.T) {
	tt := []struct {
		name string
		tags trace.RPCTags
	}{
		{
			name: "all tags",
			tags: trace.RPCTags{
				Kind:         ext.SpanKindRPCClientEnum,
				PeerAddr:     "http://internal.service.io/",
				PeerHostname: "internal.service.io",
				PeerIPv4:     net.IPv4(127, 0, 0, 1),
				PeerIPv6:     net.IPv6zero,
				PeerPort:     8080,
				PeerService:  "service",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test").(*mocktracer.MockSpan)
			trace.SetRPCTags(span, tc.tags)
			span.Finish()
			ensureRPCTagsSet(t, tc.tags, span.Tags())
		})
	}
}

func ensureRPCTagsSet(t *testing.T, rpcTags trace.RPCTags, tags map[string]interface{}) {
	key := string(ext.SpanKind)
	kind, ok := tags[key]
	switch {
	case rpcTags.Kind == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, kind)
	case rpcTags.Kind != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case rpcTags.Kind != "" && ok:
		var k string
		switch v := kind.(type) {
		case string:
			k = v
		case ext.SpanKindEnum:
			k = string(v)
		}
		if got, want := k, string(rpcTags.Kind); got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}
	ensurePeerTagsSet(
		t,
		tags,
		rpcTags.PeerAddr,
		rpcTags.PeerHostname,
		rpcTags.PeerIPv4.String(),
		rpcTags.PeerIPv6.String(),
		rpcTags.PeerService,
		rpcTags.PeerPort,
	)
}

func TestDBTags_Apply(t *testing.T) {
	tt := []struct {
		name string
		tags trace.DBTags
	}{
		{
			name: "all tags",
			tags: trace.DBTags{
				Type:         "sql",
				Instance:     "test",
				User:         "test",
				Statement:    "SELECT * FROM test",
				PeerAddr:     "http://internal.service.io/",
				PeerHostname: "internal.service.io",
				PeerIPv4:     net.IPv4(127, 0, 0, 1),
				PeerIPv6:     net.IPv6zero,
				PeerPort:     8080,
				PeerService:  "service",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test", tc.tags).(*mocktracer.MockSpan)
			span.Finish()
			ensureDBTagsSet(t, tc.tags, span.Tags())
		})
	}
}

func TestSetDBTags(t *testing.T) {
	tt := []struct {
		name string
		tags trace.DBTags
	}{
		{
			name: "all tags",
			tags: trace.DBTags{
				Type:         "sql",
				Instance:     "test",
				User:         "test",
				Statement:    "SELECT * FROM test",
				PeerAddr:     "http://internal.service.io/",
				PeerHostname: "internal.service.io",
				PeerIPv4:     net.IPv4(127, 0, 0, 1),
				PeerIPv6:     net.IPv6zero,
				PeerPort:     8080,
				PeerService:  "service",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test").(*mocktracer.MockSpan)
			trace.SetDBTags(span, tc.tags)
			span.Finish()
			ensureDBTagsSet(t, tc.tags, span.Tags())
		})
	}
}

func ensureDBTagsSet(t *testing.T, dbTags trace.DBTags, tags map[string]interface{}) {
	key := string(ext.DBType)
	dbType, ok := tags[key]
	switch {
	case dbTags.Type == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, dbType)
	case dbTags.Type != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case dbTags.Type != "" && ok:
		if got, want := dbType, dbTags.Type; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.DBInstance)
	dbInstance, ok := tags[key]
	switch {
	case dbTags.Instance == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, dbInstance)
	case dbTags.Instance != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case dbTags.Instance != "" && ok:
		if got, want := dbInstance, dbTags.Instance; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.DBUser)
	dbUser, ok := tags[key]
	switch {
	case dbTags.User == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, dbUser)
	case dbTags.User != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case dbTags.User != "" && ok:
		if got, want := dbUser, dbTags.User; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.DBStatement)
	dbStmt, ok := tags[key]
	switch {
	case dbTags.Statement == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, dbStmt)
	case dbTags.Statement != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case dbTags.Statement != "" && ok:
		if got, want := dbStmt, dbTags.Statement; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	ensurePeerTagsSet(
		t,
		tags,
		dbTags.PeerAddr,
		dbTags.PeerHostname,
		dbTags.PeerIPv4.String(),
		dbTags.PeerIPv6.String(),
		dbTags.PeerService,
		dbTags.PeerPort,
	)
}

func TestHTTPTags_Apply(t *testing.T) {
	tt := []struct {
		name string
		tags trace.HTTPTags
	}{
		{
			name: "all tags",
			tags: trace.HTTPTags{
				Method:     http.MethodGet,
				URL:        "http://example.com/",
				StatusCode: http.StatusOK,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test", tc.tags).(*mocktracer.MockSpan)
			span.Finish()
			ensureHTTPTagsSet(t, tc.tags, span.Tags())
		})
	}
}

func TestSetHTTPTags(t *testing.T) {
	tt := []struct {
		name string
		tags trace.HTTPTags
	}{
		{
			name: "all tags",
			tags: trace.HTTPTags{
				Method:     http.MethodGet,
				URL:        "http://example.com/",
				StatusCode: http.StatusOK,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test").(*mocktracer.MockSpan)
			trace.SetHTTPTags(span, tc.tags)
			span.Finish()
			ensureHTTPTagsSet(t, tc.tags, span.Tags())
		})
	}
}

func ensureHTTPTagsSet(t *testing.T, httpTags trace.HTTPTags, tags map[string]interface{}) {
	key := string(ext.HTTPMethod)
	httpMethod, ok := tags[key]
	switch {
	case httpTags.Method == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, httpMethod)
	case httpTags.Method != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case httpTags.Method != "" && ok:
		if got, want := httpMethod, httpTags.Method; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.HTTPUrl)
	httpURL, ok := tags[key]
	switch {
	case httpTags.URL == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, httpURL)
	case httpTags.URL != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case httpTags.URL != "" && ok:
		if got, want := httpURL, httpTags.URL; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.HTTPStatusCode)
	httpStatusCode, ok := tags[key]
	switch {
	case httpTags.StatusCode == 0 && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, httpStatusCode)
	case httpTags.StatusCode != 0 && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case httpTags.StatusCode != 0 && ok:
		var code uint16
		switch v := httpStatusCode.(type) {
		case int:
			code = uint16(v)
		case uint16:
			code = v
		}
		if got, want := code, uint16(httpTags.StatusCode); got != want {
			t.Errorf("tag %q: got %d, want %d\n", key, got, want)
		}
	}
}

func ensurePeerTagsSet(t *testing.T, tags map[string]interface{}, addr, hostname, ipv4, ipv6, service string, port uint16) {
	key := string(ext.PeerAddress)
	peerAddr, ok := tags[key]
	switch {
	case addr == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, peerAddr)
	case addr != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case addr != "" && ok:
		if got, want := peerAddr, string(addr); got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.PeerHostname)
	peerHost, ok := tags[key]
	switch {
	case hostname == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, peerHost)
	case hostname != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case hostname != "" && ok:
		if got, want := peerHost, string(hostname); got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.PeerHostIPv4)
	peerIPv4, ok := tags[key]
	switch {
	case ipv4 == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, peerIPv4)
	case ipv4 != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case ipv4 != "" && ok:
		if got, want := peerIPv4, ipv4; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.PeerHostIPv6)
	peerIPv6, ok := tags[key]
	switch {
	case ipv6 == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, peerIPv6)
	case ipv6 != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case ipv6 != "" && ok:
		if got, want := peerIPv6, ipv6; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.PeerPort)
	peerPort, ok := tags[key]
	switch {
	case port == 0 && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, peerPort)
	case port != 0 && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case port != 0 && ok:
		if got, want := peerPort, port; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}

	key = string(ext.PeerService)
	peerService, ok := tags[key]
	switch {
	case service == "" && ok:
		t.Errorf("tag %q: unexpected value %q\n", key, peerService)
	case service != "" && !ok:
		t.Errorf("tag %q: expected value\n", key)
	case service != "" && ok:
		if got, want := peerService, service; got != want {
			t.Errorf("tag %q: got %q, want %q\n", key, got, want)
		}
	}
}
