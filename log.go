package trace

import (
	"encoding/json"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// Standard opentracing log field names.
// See https://github.com/opentracing/specification/blob/master/semantic_conventions.md#log-fields-table.
const (
	LogFieldErrorKind   = "error.kind"
	LogFieldErrorObject = "error.object"
	LogFieldEvent       = "event"
	LogFieldMessage     = "message"
	LogFieldStack       = "stack"

	LogEventError = "error"
)

// LogFields is a map of opentracing span log field names to values.
type LogFields map[string]interface{}

// Encode returns a new map with the values for each key JSON marshaled. If the
// value for a key could not be marshaled, the original value is preserved.
func (f LogFields) Encode() map[string]interface{} {
	encoded := make(map[string]interface{})
	for k, v := range f {
		if b, err := json.Marshal(v); err == nil {
			encoded[k] = string(b)
		} else {
			encoded[k] = v
		}
	}
	return encoded
}

// LogError logs an error for an opentracing span, setting the standard error
// tags and log fields.
func LogError(span opentracing.Span, err error) {
	if span == nil || err == nil {
		return
	}
	ext.Error.Set(span, true)
	span.LogFields(
		log.String(LogFieldEvent, LogEventError),
		log.String(LogFieldErrorKind, fmt.Sprintf("%T", errors.Cause(err))),
		log.String(LogFieldMessage, err.Error()),
	)
}

// LogErrorf logs an error with the specified format for an opentracing span,
// setting the standard error tags and log fields.
func LogErrorf(span opentracing.Span, err error, format string, args ...interface{}) {
	if span == nil || err == nil {
		return
	}
	ext.Error.Set(span, true)
	span.LogFields(
		log.String(LogFieldEvent, LogEventError),
		log.String(LogFieldErrorKind, fmt.Sprintf("%T", errors.Cause(err))),
		log.String(LogFieldMessage, errors.Wrapf(err, format, args...).Error()),
	)
}

// LogErrorWithFields logs an error with the specified extra lof fields for an
// opentracing span, setting the standard error tags and log fields. The log
// field names "event", "error.kind", and "message" are reserved and will be
// ignored if set in the specified fields.
func LogErrorWithFields(span opentracing.Span, err error, fields map[string]interface{}) {
	if span == nil || err == nil {
		return
	}
	ext.Error.Set(span, true)
	kvs := []interface{}{
		LogFieldEvent, LogEventError,
		LogFieldErrorKind, fmt.Sprintf("%T", errors.Cause(err)),
		LogFieldMessage, err.Error(),
	}
	for k, v := range fields {
		if k == LogFieldEvent || k == LogFieldErrorKind || k == LogFieldMessage {
			continue
		}
		kvs = append(kvs, k, v)
	}
	span.LogKV(kvs...)
}
