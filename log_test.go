package trace_test

import (
	"fmt"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/pkg/errors"

	"github.com/code-willing/trace"
)

func init() {
	opentracing.SetGlobalTracer(mocktracer.New())
}

func TestLogError(t *testing.T) {
	tt := []struct {
		name string
		err  error
	}{
		{
			name: "error",
			err:  errors.New("error"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test").(*mocktracer.MockSpan)
			trace.LogError(span, tc.err)
			span.Finish()

			if span.Tag(string(ext.Error)) == nil {
				t.Error("expected error tag")
			}
			logs := span.Logs()
			if got, want := len(logs), 1; got != want {
				t.Fatalf("logs: got %d, want %d", got, want)
			}
			logFields := logs[0].Fields
			if got, want := len(logFields), 3; got != want {
				t.Fatalf("log fields: got %d, want %d", got, want)
			}
			for _, field := range logFields {
				switch field.Key {
				case trace.LogFieldEvent:
					if got, want := field.ValueString, trace.LogEventError; got != want {
						t.Errorf("log field: event: got %q, want %q\n", got, want)
					}
				case trace.LogFieldErrorKind:
					if got, want := field.ValueString, fmt.Sprintf("%T", errors.Cause(tc.err)); got != want {
						t.Errorf("log field: error.kind: got %q, want %q\n", got, want)
					}
				case trace.LogFieldMessage:
					if got, want := field.ValueString, tc.err.Error(); got != want {
						t.Errorf("log field: message: got %q, want %q\n", got, want)
					}
				}
			}
		})
	}
}

func TestLogErrorf(t *testing.T) {
	tt := []struct {
		name   string
		err    error
		format string
		args   []interface{}
	}{
		{
			name:   "error",
			err:    errors.New("error"),
			format: "foo: %s",
			args:   []interface{}{"bar"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test").(*mocktracer.MockSpan)
			trace.LogErrorf(span, tc.err, tc.format, tc.args...)
			span.Finish()

			if span.Tag(string(ext.Error)) == nil {
				t.Error("expected error tag")
			}
			logs := span.Logs()
			if got, want := len(logs), 1; got != want {
				t.Fatalf("logs: got %d, want %d", got, want)
			}
			logFields := logs[0].Fields
			if got, want := len(logFields), 3; got != want {
				t.Fatalf("log fields: got %d, want %d", got, want)
			}
			for _, field := range logFields {
				switch field.Key {
				case trace.LogFieldEvent:
					if got, want := field.ValueString, trace.LogEventError; got != want {
						t.Errorf("log field: event: got %q, want %q\n", got, want)
					}
				case trace.LogFieldErrorKind:
					if got, want := field.ValueString, fmt.Sprintf("%T", errors.Cause(tc.err)); got != want {
						t.Errorf("log field: error.kind: got %q, want %q\n", got, want)
					}
				case trace.LogFieldMessage:
					errMsg := errors.Wrapf(tc.err, tc.format, tc.args...).Error()
					if got, want := field.ValueString, errMsg; got != want {
						t.Errorf("log field: message: got %q, want %q\n", got, want)
					}
				}
			}
		})
	}
}

func TestLogErrorWithFields(t *testing.T) {
	tt := []struct {
		name   string
		err    error
		fields map[string]interface{}
	}{
		{
			name: "error",
			err:  errors.New("error"),
			fields: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			span := opentracing.StartSpan("test").(*mocktracer.MockSpan)
			trace.LogErrorWithFields(span, tc.err, tc.fields)
			span.Finish()

			if span.Tag(string(ext.Error)) == nil {
				t.Error("expected error tag")
			}
			logs := span.Logs()
			if got, want := len(logs), 1; got != want {
				t.Fatalf("logs: got %d, want %d", got, want)
			}
			logFields := logs[0].Fields
			if got, want := len(logFields), 3+len(tc.fields); got != want {
				t.Fatalf("log fields: got %d, want %d", got, want)
			}
			for _, field := range logFields {
				switch field.Key {
				case trace.LogFieldEvent:
					if got, want := field.ValueString, trace.LogEventError; got != want {
						t.Errorf("log field: event: got %q, want %q\n", got, want)
					}
				case trace.LogFieldErrorKind:
					if got, want := field.ValueString, fmt.Sprintf("%T", errors.Cause(tc.err)); got != want {
						t.Errorf("log field: error.kind: got %q, want %q\n", got, want)
					}
				case trace.LogFieldMessage:
					if got, want := field.ValueString, tc.err.Error(); got != want {
						t.Errorf("log field: message: got %q, want %q\n", got, want)
					}
				default:
					v, ok := tc.fields[field.Key]
					if !ok {
						t.Errorf("log field: %s: expected value\n", field.Key)
					}
					if got, want := field.ValueString, v; got != want {
						t.Errorf("log field: %s: got %q, want %q\n", field.Key, got, want)
					}
				}
			}
		})
	}
}
