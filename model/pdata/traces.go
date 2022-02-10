// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pdata // import "go.opentelemetry.io/collector/model/pdata"

import "go.opentelemetry.io/collector/model/internal/pdata"

// TracesMarshaler marshals pdata.Traces into bytes.
type TracesMarshaler = pdata.TracesMarshaler

// TracesUnmarshaler unmarshalls bytes into pdata.Traces.
type TracesUnmarshaler = pdata.TracesUnmarshaler

// TracesSizer is an optional interface implemented by the TracesMarshaler,
// that calculates the size of a marshaled Traces.
type TracesSizer = pdata.TracesSizer

// Traces is the top-level struct that is propagated through the traces pipeline.
type Traces = pdata.Traces

// NewTraces creates a new Traces.
var NewTraces = pdata.NewTraces

// TracesFromInternalRep creates Traces from the internal representation.
// Should not be used outside this module.
// TODO: Can be removed, it's used in the model package only.
var TracesFromInternalRep = pdata.TracesFromInternalRep

// TraceState is a string representing the tracestate in w3c-trace-context format: https://www.w3.org/TR/trace-context/#tracestate-header
type TraceState = pdata.TraceState

const (
	// TraceStateEmpty represents the empty TraceState.
	TraceStateEmpty = pdata.TraceStateEmpty
)

// SpanKind is the type of span. Can be used to specify additional relationships between spans
// in addition to a parent/child relationship.
type SpanKind = pdata.SpanKind

const (
	// SpanKindUnspecified represents that the SpanKind is unspecified, it MUST NOT be used.
	SpanKindUnspecified = pdata.SpanKindUnspecified
	// SpanKindInternal indicates that the span represents an internal operation within an application,
	// as opposed to an operation happening at the boundaries. Default value.
	SpanKindInternal = pdata.SpanKindInternal
	// SpanKindServer indicates that the span covers server-side handling of an RPC or other
	// remote network request.
	SpanKindServer = pdata.SpanKindServer
	// SpanKindClient indicates that the span describes a request to some remote service.
	SpanKindClient = pdata.SpanKindClient
	// SpanKindProducer indicates that the span describes a producer sending a message to a broker.
	// Unlike CLIENT and SERVER, there is often no direct critical path latency relationship
	// between producer and consumer spans.
	// A PRODUCER span ends when the message was accepted by the broker while the logical processing of
	// the message might span a much longer time.
	SpanKindProducer = pdata.SpanKindProducer
	// SpanKindConsumer indicates that the span describes consumer receiving a message from a broker.
	// Like the PRODUCER kind, there is often no direct critical path latency relationship between
	// producer and consumer spans.
	SpanKindConsumer = pdata.SpanKindConsumer
)

// StatusCode mirrors the codes defined at
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/trace/api.md#set-status
type StatusCode = pdata.StatusCode

const (
	StatusCodeUnset = pdata.StatusCodeUnset
	StatusCodeOk    = pdata.StatusCodeOk
	StatusCodeError = pdata.StatusCodeError
)
