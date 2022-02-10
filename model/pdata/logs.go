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

// LogsMarshaler marshals pdata.Logs into bytes.
type LogsMarshaler = pdata.LogsMarshaler

// LogsUnmarshaler unmarshalls bytes into pdata.Logs.
type LogsUnmarshaler = pdata.LogsUnmarshaler

// LogsSizer is an optional interface implemented by the LogsMarshaler,
// that calculates the size of a marshaled Logs.
type LogsSizer = pdata.LogsSizer

// Logs is the top-level struct that is propagated through the logs pipeline.
//
// This is a reference type (like builtin map).
//
// Must use NewLogs functions to create new instances.
// Important: zero-initialized instance is not valid for use.
type Logs = pdata.Logs

// NewLogs creates a new Logs.
var NewLogs = pdata.NewLogs

// LogsFromInternalRep creates the internal Logs representation from the ProtoBuf. Should
// not be used outside this module. This is intended to be used only by OTLP exporter and
// File exporter, which legitimately need to work with OTLP Protobuf structs.
// TODO: Can be removed, it's used in the model package only.
var LogsFromInternalRep = pdata.LogsFromInternalRep

// Deprecated: use LogRecordSlice
type LogSlice = LogRecordSlice

// Deprecated: use NewLogRecordSlice
var NewLogSlice = NewLogRecordSlice

// SeverityNumber is the public alias of otlplogs.SeverityNumber from internal package.
type SeverityNumber = pdata.SeverityNumber

const (
	SeverityNumberUNDEFINED = pdata.SeverityNumberUNDEFINED
	SeverityNumberTRACE     = pdata.SeverityNumberTRACE
	SeverityNumberTRACE2    = pdata.SeverityNumberTRACE2
	SeverityNumberTRACE3    = pdata.SeverityNumberTRACE3
	SeverityNumberTRACE4    = pdata.SeverityNumberTRACE4
	SeverityNumberDEBUG     = pdata.SeverityNumberDEBUG
	SeverityNumberDEBUG2    = pdata.SeverityNumberDEBUG2
	SeverityNumberDEBUG3    = pdata.SeverityNumberDEBUG3
	SeverityNumberDEBUG4    = pdata.SeverityNumberDEBUG4
	SeverityNumberINFO      = pdata.SeverityNumberINFO
	SeverityNumberINFO2     = pdata.SeverityNumberINFO2
	SeverityNumberINFO3     = pdata.SeverityNumberINFO3
	SeverityNumberINFO4     = pdata.SeverityNumberINFO4
	SeverityNumberWARN      = pdata.SeverityNumberWARN
	SeverityNumberWARN2     = pdata.SeverityNumberWARN2
	SeverityNumberWARN3     = pdata.SeverityNumberWARN3
	SeverityNumberWARN4     = pdata.SeverityNumberWARN4
	SeverityNumberERROR     = pdata.SeverityNumberERROR
	SeverityNumberERROR2    = pdata.SeverityNumberERROR2
	SeverityNumberERROR3    = pdata.SeverityNumberERROR3
	SeverityNumberERROR4    = pdata.SeverityNumberERROR4
	SeverityNumberFATAL     = pdata.SeverityNumberFATAL
	SeverityNumberFATAL2    = pdata.SeverityNumberFATAL2
	SeverityNumberFATAL3    = pdata.SeverityNumberFATAL3
	SeverityNumberFATAL4    = pdata.SeverityNumberFATAL4
)
