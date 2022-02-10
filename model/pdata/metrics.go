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

// MetricsMarshaler marshals pdata.Metrics into bytes.
type MetricsMarshaler = pdata.MetricsMarshaler

// MetricsUnmarshaler unmarshalls bytes into pdata.Metrics.
type MetricsUnmarshaler = pdata.MetricsUnmarshaler

// MetricsSizer is an optional interface implemented by the MetricsMarshaler,
// that calculates the size of a marshaled Metrics.
type MetricsSizer = pdata.MetricsSizer

// Metrics is an opaque interface that allows transition to the new internal Metrics data, but also facilitates the
// transition to the new components, especially for traces.
//
// Outside of the core repository, the metrics pipeline cannot be converted to the new model since data.MetricData is
// part of the internal package.
type Metrics = pdata.Metrics

// NewMetrics creates a new Metrics.
var NewMetrics = pdata.NewMetrics

// MetricsFromInternalRep creates Metrics from the internal representation.
// Should not be used outside this module.
// TODO: Can be removed, it's used in the model package only.
var MetricsFromInternalRep = pdata.MetricsFromInternalRep

// MetricDataType specifies the type of data in a Metric.
type MetricDataType = pdata.MetricDataType

const (
	MetricDataTypeNone                 = pdata.MetricDataTypeNone
	MetricDataTypeGauge                = pdata.MetricDataTypeGauge
	MetricDataTypeSum                  = pdata.MetricDataTypeSum
	MetricDataTypeHistogram            = pdata.MetricDataTypeHistogram
	MetricDataTypeExponentialHistogram = pdata.MetricDataTypeExponentialHistogram
	MetricDataTypeSummary              = pdata.MetricDataTypeSummary
)

// MetricAggregationTemporality defines how a metric aggregator reports aggregated values.
// It describes how those values relate to the time interval over which they are aggregated.
type MetricAggregationTemporality = pdata.MetricAggregationTemporality

const (
	// MetricAggregationTemporalityUnspecified is the default MetricAggregationTemporality, it MUST NOT be used.
	MetricAggregationTemporalityUnspecified = pdata.MetricAggregationTemporalityUnspecified
	// MetricAggregationTemporalityDelta is a MetricAggregationTemporality for a metric aggregator which reports changes since last report time.
	MetricAggregationTemporalityDelta = pdata.MetricAggregationTemporalityDelta
	// MetricAggregationTemporalityCumulative is a MetricAggregationTemporality for a metric aggregator which reports changes since a fixed start time.
	MetricAggregationTemporalityCumulative = pdata.MetricAggregationTemporalityCumulative
)

// MetricDataPointFlags defines how a metric aggregator reports aggregated values.
// It describes how those values relate to the time interval over which they are aggregated.
type MetricDataPointFlags = pdata.MetricDataPointFlags

const (
	// MetricDataPointFlagsNone is the default MetricDataPointFlags.
	MetricDataPointFlagsNone = pdata.MetricDataPointFlagsNone
)

// NewMetricDataPointFlags returns a new MetricDataPointFlags combining the flags passed
// in as parameters.
var NewMetricDataPointFlags = pdata.NewMetricDataPointFlags

// MetricDataPointFlag allow users to configure DataPointFlags. This is achieved via NewMetricDataPointFlags.
// The separation between MetricDataPointFlags and MetricDataPointFlag exists to prevent users accidentally
// comparing the value of individual flags with MetricDataPointFlags. Instead, users must use the HasFlag method.
type MetricDataPointFlag = pdata.MetricDataPointFlag

const (
	// MetricDataPointFlagNoRecordedValue is flag for a metric aggregator which reports changes since last report time.
	MetricDataPointFlagNoRecordedValue = pdata.MetricDataPointFlagNoRecordedValue
)

// MetricValueType specifies the type of NumberDataPoint.
type MetricValueType = pdata.MetricValueType

const (
	MetricValueTypeNone   = pdata.MetricValueTypeNone
	MetricValueTypeInt    = pdata.MetricValueTypeInt
	MetricValueTypeDouble = pdata.MetricValueTypeDouble
)
