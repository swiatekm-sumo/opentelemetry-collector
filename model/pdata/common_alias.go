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

// This file contains aliases to data structures that are common for all
// signal types, such as timestamps, attributes, etc.

import "go.opentelemetry.io/collector/model/internal/pdata"

// AnyValueType is an alias for pdata.AnyValueType type.
type AnyValueType = pdata.AnyValueType

// AttributeValueType is an alias for pdata.AnyValueType type.
// Deprecated: [v0.47.0] Use AnyValueType instead.
type AttributeValueType = pdata.AnyValueType

const (
	AnyValueTypeEmpty  = pdata.AnyValueTypeEmpty
	AnyValueTypeString = pdata.AnyValueTypeString
	AnyValueTypeInt    = pdata.AnyValueTypeInt
	AnyValueTypeDouble = pdata.AnyValueTypeDouble
	AnyValueTypeBool   = pdata.AnyValueTypeBool
	AnyValueTypeMap    = pdata.AnyValueTypeMap
	AnyValueTypeArray  = pdata.AnyValueTypeArray
	AnyValueTypeBytes  = pdata.AnyValueTypeBytes

	// Deprecated: [v0.47.0] Use AnyValueTypeEmpty instead.
	AttributeValueTypeEmpty = pdata.AnyValueTypeEmpty

	// Deprecated: [v0.47.0] Use AnyValueTypeString instead.
	AttributeValueTypeString = pdata.AnyValueTypeString

	// Deprecated: [v0.47.0] Use AnyValueTypeInt instead.
	AttributeValueTypeInt = pdata.AnyValueTypeInt

	// Deprecated: [v0.47.0] Use AnyValueTypeDouble instead.
	AttributeValueTypeDouble = pdata.AnyValueTypeDouble

	// Deprecated: [v0.47.0] Use AnyValueTypeBool instead.
	AttributeValueTypeBool = pdata.AnyValueTypeBool

	// Deprecated: [v0.47.0] Use AnyValueTypeMap instead.
	AttributeValueTypeMap = pdata.AnyValueTypeMap

	// Deprecated: [v0.47.0] Use AnyValueTypeArray instead.
	AttributeValueTypeArray = pdata.AnyValueTypeArray

	// Deprecated: [v0.47.0] Use AnyValueTypeBytes instead.
	AttributeValueTypeBytes = pdata.AnyValueTypeBytes
)

// AnyValue is an alias for pdata.AnyValue struct.
type AnyValue = pdata.AnyValue

// AttributeValue is an alias for pdata.AnyValue struct.
// Deprecated: [v0.47.0] Use AnyValue instead.
type AttributeValue = pdata.AnyValue

// Aliases for functions to create pdata.AnyValue.
var (
	NewAnyValueEmpty  = pdata.NewAnyValueEmpty
	NewAnyValueString = pdata.NewAnyValueString
	NewAnyValueInt    = pdata.NewAnyValueInt
	NewAnyValueDouble = pdata.NewAnyValueDouble
	NewAnyValueBool   = pdata.NewAnyValueBool
	NewAnyValueMap    = pdata.NewAnyValueMap
	NewAnyValueArray  = pdata.NewAnyValueArray
	NewAnyValueBytes  = pdata.NewAnyValueBytes

	// Deprecated: [v0.47.0] Use NewAnyValueEmpty instead.
	NewAttributeValueEmpty = pdata.NewAnyValueEmpty

	// Deprecated: [v0.47.0] Use NewAnyValueString instead.
	NewAttributeValueString = pdata.NewAnyValueString

	// Deprecated: [v0.47.0] Use NewAnyValueInt instead.
	NewAttributeValueInt = pdata.NewAnyValueInt

	// Deprecated: [v0.47.0] Use NewAnyValueDouble instead.
	NewAttributeValueDouble = pdata.NewAnyValueDouble

	// Deprecated: [v0.47.0] Use NewAnyValueBool instead.
	NewAttributeValueBool = pdata.NewAnyValueBool

	// Deprecated: [v0.47.0] Use NewAnyValueMap instead.
	NewAttributeValueMap = pdata.NewAnyValueMap

	// Deprecated: [v0.47.0] Use NewAnyValueArray instead.
	NewAttributeValueArray = pdata.NewAnyValueArray

	// Deprecated: [v0.47.0] Use NewAnyValueBytes instead.
	NewAttributeValueBytes = pdata.NewAnyValueBytes
)

// AttributeValueSlice is an alias for pdata.AttributeValueSlice struct.
// Deprecated: [v0.47.0] Use AttributeValueSlice instead.
type AttributeValueSlice = pdata.AttributeValueSlice

// NewAttributeValueSlice is an alias for a function to create AttributeValueSlice.
// Deprecated: [v0.47.0] Use NewAttributeValueSlice instead.
var NewAttributeValueSlice = pdata.NewAttributeValueSlice

// AttributeMap is an alias for pdata.AttributeMap struct.
type AttributeMap = pdata.AttributeMap

// Aliases for functions to create pdata.AttributeMap.
var (
	NewAttributeMap        = pdata.NewAttributeMap
	NewAttributeMapFromMap = pdata.NewAttributeMapFromMap
)
