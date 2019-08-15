package trace

import (
	"net"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Ensure RPCTags implements the opentracing.StartSpanOption interface.
var _ opentracing.StartSpanOption = (*RPCTags)(nil)

// RPCTags is an opentracing.StartSpanOption that sets the standard RPC tags.
//
// See https://github.com/opentracing/specification/blob/master/semantic_conventions.md#rpcs.
type RPCTags struct {
	Kind ext.SpanKindEnum // The span kind, "client" or "server".

	// Optional tags that describe the RPC peer.
	PeerAddr     string // The remote address.
	PeerHostname string // The remote hostname.
	PeerIPv4     net.IP // The remote IPv4 address.
	PeerIPv6     net.IP // The remote IPv6 address.
	PeerPort     uint16 // The remote port.
	PeerService  string // The remote service name.
}

// Apply implements the opentracing.StartSpanOption interface.
func (t RPCTags) Apply(opts *opentracing.StartSpanOptions) {
	if opts == nil {
		return
	}
	if opts.Tags == nil {
		opts.Tags = make(map[string]interface{})
	}
	if t.Kind == ext.SpanKindRPCClientEnum || t.Kind == ext.SpanKindRPCServerEnum {
		opts.Tags[string(ext.SpanKind)] = string(t.Kind)
	}
	applyPeerTags(
		opts,
		t.PeerAddr,
		t.PeerHostname,
		t.PeerService,
		t.PeerIPv4,
		t.PeerIPv6,
		t.PeerPort,
	)
}

// SetRPCTags sets the standard RPC tags on the specified span.
func SetRPCTags(span opentracing.Span, t RPCTags) {
	if span == nil {
		return
	}
	if t.Kind == ext.SpanKindRPCClientEnum || t.Kind == ext.SpanKindRPCServerEnum {
		ext.SpanKind.Set(span, t.Kind)
	}
	setPeerTags(
		span,
		t.PeerAddr,
		t.PeerHostname,
		t.PeerService,
		t.PeerIPv4,
		t.PeerIPv6,
		t.PeerPort,
	)
}

// Ensure DBTags implements the opentracing.StartSpanOption interface.
var _ opentracing.StartSpanOption = (*DBTags)(nil)

// DBTags is an opentracing.StartSpanOption that sets the standard database tags
// for a database client call.
//
// See https://github.com/opentracing/specification/blob/master/semantic_conventions.md#database-client-calls.
type DBTags struct {
	Type      string // The database type.
	Instance  string // The database instance name.
	User      string // The username of the database accessor.
	Statement string // The database statement used.

	// Optional tags that describe the database peer.
	PeerAddr     string // The remote address.
	PeerHostname string // The remote hostname.
	PeerIPv4     net.IP // The remote IPv4 address.
	PeerIPv6     net.IP // The remote IPv6 address.
	PeerPort     uint16 // The remote port.
	PeerService  string // The remote service name.
}

// Apply implements the opentracing.StartSpanOption interface.
func (t DBTags) Apply(opts *opentracing.StartSpanOptions) {
	if opts == nil {
		return
	}
	if opts.Tags == nil {
		opts.Tags = make(map[string]interface{})
	}
	opts.Tags[string(ext.SpanKind)] = ext.SpanKindRPCClientEnum
	if t.Type != "" {
		opts.Tags[string(ext.DBType)] = strings.ToLower(t.Type)
	}
	if t.Instance != "" {
		opts.Tags[string(ext.DBInstance)] = t.Instance
	}
	if t.User != "" {
		opts.Tags[string(ext.DBUser)] = t.User
	}
	if t.Statement != "" {
		opts.Tags[string(ext.DBStatement)] = t.Statement
	}
	applyPeerTags(
		opts,
		t.PeerAddr,
		t.PeerHostname,
		t.PeerService,
		t.PeerIPv4,
		t.PeerIPv6,
		t.PeerPort,
	)
}

// SetDBTags sets the standard database tags on the specified span.
func SetDBTags(span opentracing.Span, t DBTags) {
	if span == nil {
		return
	}
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	if t.Type != "" {
		ext.DBType.Set(span, strings.ToLower(t.Type))
	}
	if t.Instance != "" {
		ext.DBInstance.Set(span, t.Instance)
	}
	if t.User != "" {
		ext.DBUser.Set(span, t.User)
	}
	if t.Statement != "" {
		ext.DBStatement.Set(span, t.Statement)
	}
	setPeerTags(
		span,
		t.PeerAddr,
		t.PeerHostname,
		t.PeerService,
		t.PeerIPv4,
		t.PeerIPv6,
		t.PeerPort,
	)
}

// Ensure HTTPTags implements the opentracing.StartSpanOption interface.
var _ opentracing.StartSpanOption = (*HTTPTags)(nil)

// HTTPTags is an opentracing.StartSpanOption that sets the standard HTTP tags.
//
// See https://github.com/opentracing/specification/blob/master/semantic_conventions.md#span-tags-table.
type HTTPTags struct {
	Method     string // The HTTP request method.
	URL        string // The HTTP request URL.
	StatusCode int    // The HTTP response status code.
}

// Apply implements the opentracing.StartSpanOption interface.
func (t HTTPTags) Apply(opts *opentracing.StartSpanOptions) {
	if opts == nil {
		return
	}
	if opts.Tags == nil {
		opts.Tags = make(map[string]interface{})
	}
	if t.Method != "" {
		opts.Tags[string(ext.HTTPMethod)] = t.Method
	}
	if t.URL != "" {
		opts.Tags[string(ext.HTTPUrl)] = t.URL
	}
	if t.StatusCode > 0 {
		opts.Tags[string(ext.HTTPStatusCode)] = t.StatusCode
	}
}

// SetHTTPTags sets the standard HTTP tags on the specified span.
func SetHTTPTags(span opentracing.Span, t HTTPTags) {
	if span == nil {
		return
	}
	if t.Method != "" {
		ext.HTTPMethod.Set(span, t.Method)
	}
	if t.URL != "" {
		ext.HTTPUrl.Set(span, t.URL)
	}
	if t.StatusCode > 0 {
		ext.HTTPStatusCode.Set(span, uint16(t.StatusCode))
	}
}

func applyPeerTags(opts *opentracing.StartSpanOptions, addr, hostname, service string, ipv4, ipv6 net.IP, port uint16) {
	if opts == nil {
		return
	}
	if opts.Tags == nil {
		opts.Tags = make(map[string]interface{})
	}
	if addr != "" {
		opts.Tags[string(ext.PeerAddress)] = addr
	}
	if hostname != "" {
		opts.Tags[string(ext.PeerHostname)] = hostname
	}
	if ipv4 != nil {
		opts.Tags[string(ext.PeerHostIPv4)] = ipv4.String()
	}
	if ipv6 != nil {
		opts.Tags[string(ext.PeerHostIPv6)] = ipv6.String()
	}
	if port > 0 {
		opts.Tags[string(ext.PeerPort)] = port
	}
	if service != "" {
		opts.Tags[string(ext.PeerService)] = service
	}
}

// setPeerTags sets the standard opentracing peer tags on the specified span.
func setPeerTags(span opentracing.Span, addr, hostname, service string, ipv4, ipv6 net.IP, port uint16) {
	if span == nil {
		return
	}
	if addr != "" {
		ext.PeerAddress.Set(span, addr)
	}
	if hostname != "" {
		ext.PeerHostname.Set(span, hostname)
	}
	if ipv4 != nil {
		ext.PeerHostIPv4.SetString(span, ipv4.String())
	}
	if ipv6 != nil {
		ext.PeerHostIPv6.Set(span, ipv6.String())
	}
	if port > 0 {
		ext.PeerPort.Set(span, port)
	}
	if service != "" {
		ext.PeerService.Set(span, service)
	}
}
