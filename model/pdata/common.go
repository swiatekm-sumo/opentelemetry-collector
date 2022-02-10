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

// This file contains data structures that are common for all telemetry types,
// such as timestamps, attributes, etc.

import "go.opentelemetry.io/collector/model/internal/pdata"

// AttributeValueType specifies the type of AttributeValue.
type AttributeValueType = pdata.AttributeValueType

const (
	AttributeValueTypeEmpty  = pdata.AttributeValueTypeEmpty
	AttributeValueTypeString = pdata.AttributeValueTypeString
	AttributeValueTypeInt    = pdata.AttributeValueTypeInt
	AttributeValueTypeDouble = pdata.AttributeValueTypeDouble
	AttributeValueTypeBool   = pdata.AttributeValueTypeBool
	AttributeValueTypeMap    = pdata.AttributeValueTypeMap
	AttributeValueTypeArray  = pdata.AttributeValueTypeArray
	AttributeValueTypeBytes  = pdata.AttributeValueTypeBytes
)

type AttributeValue = pdata.AttributeValue

var NewAttributeValueEmpty = pdata.NewAttributeValueEmpty

// NewAttributeValueString creates a new AttributeValue with the given string value.
var NewAttributeValueString = pdata.NewAttributeValueString

// NewAttributeValueInt creates a new AttributeValue with the given int64 value.
var NewAttributeValueInt = pdata.NewAttributeValueInt

// NewAttributeValueDouble creates a new AttributeValue with the given float64 value.
var NewAttributeValueDouble = pdata.NewAttributeValueDouble

// NewAttributeValueBool creates a new AttributeValue with the given bool value.
var NewAttributeValueBool = pdata.NewAttributeValueBool

// NewAttributeValueMap creates a new AttributeValue of map type.
var NewAttributeValueMap = pdata.NewAttributeValueMap

// NewAttributeValueArray creates a new AttributeValue of array type.
var NewAttributeValueArray = pdata.NewAttributeValueArray

// NewAttributeValueBytes creates a new AttributeValue with the given []byte value.
// The caller must ensure the []byte passed in is not modified after the call is made, sharing the data
// across multiple attributes is forbidden.
var NewAttributeValueBytes = pdata.NewAttributeValueBytes

// AttributeMap stores a map of attribute keys to values.
type AttributeMap = pdata.AttributeMap

// NewAttributeMap creates a AttributeMap with 0 elements.
var NewAttributeMap = pdata.NewAttributeMap

// NewAttributeMapFromMap creates a AttributeMap with values
// from the given map[string]AttributeValue.
var NewAttributeMapFromMap = pdata.NewAttributeMapFromMap
