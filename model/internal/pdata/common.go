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

package pdata // import "go.opentelemetry.io/collector/model/internal/pdata"

// This file contains data structures that are common for all telemetry types,
// such as timestamps, attributes, etc.

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	otlpcommon "go.opentelemetry.io/collector/model/internal/data/protogen/common/v1"
)

// AnyValueType specifies the type of AnyValue.
type AnyValueType int32

const (
	AnyValueTypeEmpty AnyValueType = iota
	AnyValueTypeString
	AnyValueTypeInt
	AnyValueTypeDouble
	AnyValueTypeBool
	AnyValueTypeMap
	AnyValueTypeArray
	AnyValueTypeBytes
)

// String returns the string representation of the AnyValueType.
func (avt AnyValueType) String() string {
	switch avt {
	case AnyValueTypeEmpty:
		return "EMPTY"
	case AnyValueTypeString:
		return "STRING"
	case AnyValueTypeBool:
		return "BOOL"
	case AnyValueTypeInt:
		return "INT"
	case AnyValueTypeDouble:
		return "DOUBLE"
	case AnyValueTypeMap:
		return "MAP"
	case AnyValueTypeArray:
		return "ARRAY"
	case AnyValueTypeBytes:
		return "BYTES"
	}
	return ""
}

// AnyValue is a mutable cell containing the value of an attribute. Typically used in AttributeMap.
// Must use one of NewAnyValue+ functions below to create new instances.
//
// Intended to be passed by value since internally it is just a pointer to actual
// value representation. For the same reason passing by value and calling setters
// will modify the original, e.g.:
//
//   func f1(val AnyValue) { val.SetIntVal(234) }
//   func f2() {
//       v := NewAnyValueString("a string")
//       f1(v)
//       _ := v.Type() // this will return AnyValueTypeInt
//   }
//
// Important: zero-initialized instance is not valid for use. All AnyValue functions below must
// be called only on instances that are created via NewAnyValue+ functions.
type AnyValue struct {
	orig *otlpcommon.AnyValue
}

func newAnyValue(orig *otlpcommon.AnyValue) AnyValue {
	return AnyValue{orig}
}

// NewAnyValueEmpty creates a new AnyValue with an empty value.
func NewAnyValueEmpty() AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{}}
}

// NewAnyValueString creates a new AnyValue with the given string value.
func NewAnyValueString(v string) AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_StringValue{StringValue: v}}}
}

// NewAnyValueInt creates a new AnyValue with the given int64 value.
func NewAnyValueInt(v int64) AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_IntValue{IntValue: v}}}
}

// NewAnyValueDouble creates a new AnyValue with the given float64 value.
func NewAnyValueDouble(v float64) AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_DoubleValue{DoubleValue: v}}}
}

// NewAnyValueBool creates a new AnyValue with the given bool value.
func NewAnyValueBool(v bool) AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_BoolValue{BoolValue: v}}}
}

// NewAnyValueMap creates a new AnyValue of map type.
func NewAnyValueMap() AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_KvlistValue{KvlistValue: &otlpcommon.KeyValueList{}}}}
}

// NewAnyValueArray creates a new AnyValue of array type.
func NewAnyValueArray() AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{}}}}
}

// NewAnyValueBytes creates a new AnyValue with the given []byte value.
// The caller must ensure the []byte passed in is not modified after the call is made, sharing the data
// across multiple attributes is forbidden.
func NewAnyValueBytes(v []byte) AnyValue {
	return AnyValue{orig: &otlpcommon.AnyValue{Value: &otlpcommon.AnyValue_BytesValue{BytesValue: v}}}
}

// Type returns the type of the value for this AnyValue.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) Type() AnyValueType {
	if a.orig.Value == nil {
		return AnyValueTypeEmpty
	}
	switch a.orig.Value.(type) {
	case *otlpcommon.AnyValue_StringValue:
		return AnyValueTypeString
	case *otlpcommon.AnyValue_BoolValue:
		return AnyValueTypeBool
	case *otlpcommon.AnyValue_IntValue:
		return AnyValueTypeInt
	case *otlpcommon.AnyValue_DoubleValue:
		return AnyValueTypeDouble
	case *otlpcommon.AnyValue_KvlistValue:
		return AnyValueTypeMap
	case *otlpcommon.AnyValue_ArrayValue:
		return AnyValueTypeArray
	case *otlpcommon.AnyValue_BytesValue:
		return AnyValueTypeBytes
	}
	return AnyValueTypeEmpty
}

// StringVal returns the string value associated with this AnyValue.
// If the Type() is not AnyValueTypeString then returns empty string.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) StringVal() string {
	return a.orig.GetStringValue()
}

// IntVal returns the int64 value associated with this AnyValue.
// If the Type() is not AnyValueTypeInt then returns int64(0).
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) IntVal() int64 {
	return a.orig.GetIntValue()
}

// DoubleVal returns the float64 value associated with this AnyValue.
// If the Type() is not AnyValueTypeDouble then returns float64(0).
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) DoubleVal() float64 {
	return a.orig.GetDoubleValue()
}

// BoolVal returns the bool value associated with this AnyValue.
// If the Type() is not AnyValueTypeBool then returns false.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) BoolVal() bool {
	return a.orig.GetBoolValue()
}

// MapVal returns the map value associated with this AnyValue.
// If the Type() is not AnyValueTypeMap then returns an empty map. Note that modifying
// such empty map has no effect on this AnyValue.
//
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) MapVal() AttributeMap {
	kvlist := a.orig.GetKvlistValue()
	if kvlist == nil {
		return NewAttributeMap()
	}
	return newAttributeMap(&kvlist.Values)
}

// SliceVal returns the slice value associated with this AnyValue.
// If the Type() is not AnyValueTypeArray then returns an empty slice. Note that modifying
// such empty slice has no effect on this AnyValue.
//
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) SliceVal() AttributeValueSlice {
	arr := a.orig.GetArrayValue()
	if arr == nil {
		return NewAttributeValueSlice()
	}
	return newAttributeValueSlice(&arr.Values)
}

// BytesVal returns the []byte value associated with this AnyValue.
// If the Type() is not AnyValueTypeBytes then returns false.
// Calling this function on zero-initialized AnyValue will cause a panic.
// Modifying the returned []byte in-place is forbidden.
func (a AnyValue) BytesVal() []byte {
	return a.orig.GetBytesValue()
}

// SetStringVal replaces the string value associated with this AnyValue,
// it also changes the type to be AnyValueTypeString.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) SetStringVal(v string) {
	a.orig.Value = &otlpcommon.AnyValue_StringValue{StringValue: v}
}

// SetIntVal replaces the int64 value associated with this AnyValue,
// it also changes the type to be AnyValueTypeInt.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) SetIntVal(v int64) {
	a.orig.Value = &otlpcommon.AnyValue_IntValue{IntValue: v}
}

// SetDoubleVal replaces the float64 value associated with this AnyValue,
// it also changes the type to be AnyValueTypeDouble.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) SetDoubleVal(v float64) {
	a.orig.Value = &otlpcommon.AnyValue_DoubleValue{DoubleValue: v}
}

// SetBoolVal replaces the bool value associated with this AnyValue,
// it also changes the type to be AnyValueTypeBool.
// Calling this function on zero-initialized AnyValue will cause a panic.
func (a AnyValue) SetBoolVal(v bool) {
	a.orig.Value = &otlpcommon.AnyValue_BoolValue{BoolValue: v}
}

// SetBytesVal replaces the []byte value associated with this AnyValue,
// it also changes the type to be AnyValueTypeBytes.
// Calling this function on zero-initialized AnyValue will cause a panic.
// The caller must ensure the []byte passed in is not modified after the call is made, sharing the data
// across multiple attributes is forbidden.
func (a AnyValue) SetBytesVal(v []byte) {
	a.orig.Value = &otlpcommon.AnyValue_BytesValue{BytesValue: v}
}

// copyTo copies the value to AnyValue. Will panic if dest is nil.
func (a AnyValue) copyTo(dest *otlpcommon.AnyValue) {
	switch v := a.orig.Value.(type) {
	case *otlpcommon.AnyValue_KvlistValue:
		kv, ok := dest.Value.(*otlpcommon.AnyValue_KvlistValue)
		if !ok {
			kv = &otlpcommon.AnyValue_KvlistValue{KvlistValue: &otlpcommon.KeyValueList{}}
			dest.Value = kv
		}
		if v.KvlistValue == nil {
			kv.KvlistValue = nil
			return
		}
		// Deep copy to dest.
		newAttributeMap(&v.KvlistValue.Values).CopyTo(newAttributeMap(&kv.KvlistValue.Values))
	case *otlpcommon.AnyValue_ArrayValue:
		av, ok := dest.Value.(*otlpcommon.AnyValue_ArrayValue)
		if !ok {
			av = &otlpcommon.AnyValue_ArrayValue{ArrayValue: &otlpcommon.ArrayValue{}}
			dest.Value = av
		}
		if v.ArrayValue == nil {
			av.ArrayValue = nil
			return
		}
		// Deep copy to dest.
		newAttributeValueSlice(&v.ArrayValue.Values).CopyTo(newAttributeValueSlice(&av.ArrayValue.Values))
	default:
		// Primitive immutable type, no need for deep copy.
		dest.Value = a.orig.Value
	}
}

// CopyTo copies the attribute to a destination.
func (a AnyValue) CopyTo(dest AnyValue) {
	a.copyTo(dest.orig)
}

// Equal checks for equality, it returns true if the objects are equal otherwise false.
func (a AnyValue) Equal(av AnyValue) bool {
	if a.orig == av.orig {
		return true
	}

	if a.orig.Value == nil || av.orig.Value == nil {
		return a.orig.Value == av.orig.Value
	}

	if a.Type() != av.Type() {
		return false
	}

	switch v := a.orig.Value.(type) {
	case *otlpcommon.AnyValue_StringValue:
		return v.StringValue == av.orig.GetStringValue()
	case *otlpcommon.AnyValue_BoolValue:
		return v.BoolValue == av.orig.GetBoolValue()
	case *otlpcommon.AnyValue_IntValue:
		return v.IntValue == av.orig.GetIntValue()
	case *otlpcommon.AnyValue_DoubleValue:
		return v.DoubleValue == av.orig.GetDoubleValue()
	case *otlpcommon.AnyValue_ArrayValue:
		vv := v.ArrayValue.GetValues()
		avv := av.orig.GetArrayValue().GetValues()
		if len(vv) != len(avv) {
			return false
		}

		for i, val := range avv {
			val := val
			newAv := newAnyValue(&vv[i])

			// According to the specification, array values must be scalar.
			if avType := newAv.Type(); avType == AnyValueTypeArray || avType == AnyValueTypeMap {
				return false
			}

			if !newAv.Equal(newAnyValue(&val)) {
				return false
			}
		}
		return true
	case *otlpcommon.AnyValue_KvlistValue:
		cc := v.KvlistValue.GetValues()
		avv := av.orig.GetKvlistValue().GetValues()
		if len(cc) != len(avv) {
			return false
		}

		am := newAttributeMap(&avv)

		for _, val := range cc {
			newAv, ok := am.Get(val.Key)
			if !ok {
				return false
			}

			if !newAv.Equal(newAnyValue(&val.Value)) {
				return false
			}
		}
		return true
	case *otlpcommon.AnyValue_BytesValue:
		return bytes.Equal(v.BytesValue, av.orig.GetBytesValue())
	}

	return false
}

// AsString converts an OTLP AnyValue object of any type to its equivalent string
// representation. This differs from StringVal which only returns a non-empty value
// if the AnyValueType is AnyValueTypeString.
func (a AnyValue) AsString() string {
	switch a.Type() {
	case AnyValueTypeEmpty:
		return ""

	case AnyValueTypeString:
		return a.StringVal()

	case AnyValueTypeBool:
		return strconv.FormatBool(a.BoolVal())

	case AnyValueTypeDouble:
		return strconv.FormatFloat(a.DoubleVal(), 'f', -1, 64)

	case AnyValueTypeInt:
		return strconv.FormatInt(a.IntVal(), 10)

	case AnyValueTypeMap:
		jsonStr, _ := json.Marshal(a.MapVal().AsRaw())
		return string(jsonStr)

	case AnyValueTypeBytes:
		return base64.StdEncoding.EncodeToString(a.BytesVal())

	case AnyValueTypeArray:
		jsonStr, _ := json.Marshal(a.SliceVal().asRaw())
		return string(jsonStr)

	default:
		return fmt.Sprintf("<Unknown OpenTelemetry attribute value type %q>", a.Type())
	}
}

func newAttributeKeyValueString(k string, v string) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	akv := AnyValue{&orig.Value}
	akv.SetStringVal(v)
	return orig
}

func newAttributeKeyValueInt(k string, v int64) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	akv := AnyValue{&orig.Value}
	akv.SetIntVal(v)
	return orig
}

func newAttributeKeyValueDouble(k string, v float64) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	akv := AnyValue{&orig.Value}
	akv.SetDoubleVal(v)
	return orig
}

func newAttributeKeyValueBool(k string, v bool) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	akv := AnyValue{&orig.Value}
	akv.SetBoolVal(v)
	return orig
}

func newAttributeKeyValueNull(k string) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	return orig
}

func newAttributeKeyValue(k string, av AnyValue) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	av.copyTo(&orig.Value)
	return orig
}

func newAttributeKeyValueBytes(k string, v []byte) otlpcommon.KeyValue {
	orig := otlpcommon.KeyValue{Key: k}
	akv := AnyValue{&orig.Value}
	akv.SetBytesVal(v)
	return orig
}

// AttributeMap stores a map of attribute keys to values.
type AttributeMap struct {
	orig *[]otlpcommon.KeyValue
}

// NewAttributeMap creates a AttributeMap with 0 elements.
func NewAttributeMap() AttributeMap {
	orig := []otlpcommon.KeyValue(nil)
	return AttributeMap{&orig}
}

// NewAttributeMapFromMap creates a AttributeMap with values
// from the given map[string]AnyValue.
func NewAttributeMapFromMap(attrMap map[string]AnyValue) AttributeMap {
	if len(attrMap) == 0 {
		kv := []otlpcommon.KeyValue(nil)
		return AttributeMap{&kv}
	}
	origs := make([]otlpcommon.KeyValue, len(attrMap))
	ix := 0
	for k, v := range attrMap {
		origs[ix].Key = k
		v.copyTo(&origs[ix].Value)
		ix++
	}
	return AttributeMap{&origs}
}

func newAttributeMap(orig *[]otlpcommon.KeyValue) AttributeMap {
	return AttributeMap{orig}
}

// Clear erases any existing entries in this AttributeMap instance.
func (am AttributeMap) Clear() {
	*am.orig = nil
}

// EnsureCapacity increases the capacity of this AttributeMap instance, if necessary,
// to ensure that it can hold at least the number of elements specified by the capacity argument.
func (am AttributeMap) EnsureCapacity(capacity int) {
	if capacity <= cap(*am.orig) {
		return
	}
	oldOrig := *am.orig
	*am.orig = make([]otlpcommon.KeyValue, 0, capacity)
	copy(*am.orig, oldOrig)
}

// Get returns the AnyValue associated with the key and true. Returned
// AnyValue is not a copy, it is a reference to the value stored in this map.
// It is allowed to modify the returned value using AnyValue.Set* functions.
// Such modification will be applied to the value stored in this map.
//
// If the key does not exist returns an invalid instance of the KeyValue and false.
// Calling any functions on the returned invalid instance will cause a panic.
func (am AttributeMap) Get(key string) (AnyValue, bool) {
	for i := range *am.orig {
		akv := &(*am.orig)[i]
		if akv.Key == key {
			return AnyValue{&akv.Value}, true
		}
	}
	return AnyValue{nil}, false
}

// Delete deletes the entry associated with the key and returns true if the key
// was present in the map, otherwise returns false.
// Deprecated: [v0.46.0] Use Remove instead.
func (am AttributeMap) Delete(key string) bool {
	return am.Remove(key)
}

// Remove removes the entry associated with the key and returns true if the key
// was present in the map, otherwise returns false.
func (am AttributeMap) Remove(key string) bool {
	for i := range *am.orig {
		akv := &(*am.orig)[i]
		if akv.Key == key {
			*akv = (*am.orig)[len(*am.orig)-1]
			*am.orig = (*am.orig)[:len(*am.orig)-1]
			return true
		}
	}
	return false
}

// RemoveIf removes the entries for which the function in question returns true
func (am AttributeMap) RemoveIf(f func(string, AnyValue) bool) {
	newLen := 0
	for i := 0; i < len(*am.orig); i++ {
		akv := &(*am.orig)[i]
		if f(akv.Key, AnyValue{&akv.Value}) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		(*am.orig)[newLen] = (*am.orig)[i]
		newLen++
	}
	*am.orig = (*am.orig)[:newLen]
}

// Insert adds the AnyValue to the map when the key does not exist.
// No action is applied to the map where the key already exists.
//
// Calling this function with a zero-initialized AnyValue struct will cause a panic.
//
// Important: this function should not be used if the caller has access to
// the raw value to avoid an extra allocation.
func (am AttributeMap) Insert(k string, v AnyValue) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValue(k, v))
	}
}

// InsertNull adds a null Value to the map when the key does not exist.
// No action is applied to the map where the key already exists.
func (am AttributeMap) InsertNull(k string) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValueNull(k))
	}
}

// InsertString adds the string Value to the map when the key does not exist.
// No action is applied to the map where the key already exists.
func (am AttributeMap) InsertString(k string, v string) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValueString(k, v))
	}
}

// InsertInt adds the int Value to the map when the key does not exist.
// No action is applied to the map where the key already exists.
func (am AttributeMap) InsertInt(k string, v int64) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValueInt(k, v))
	}
}

// InsertDouble adds the double Value to the map when the key does not exist.
// No action is applied to the map where the key already exists.
func (am AttributeMap) InsertDouble(k string, v float64) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValueDouble(k, v))
	}
}

// InsertBool adds the bool Value to the map when the key does not exist.
// No action is applied to the map where the key already exists.
func (am AttributeMap) InsertBool(k string, v bool) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValueBool(k, v))
	}
}

// InsertBytes adds the []byte Value to the map when the key does not exist.
// No action is applied to the map where the key already exists.
// The caller must ensure the []byte passed in is not modified after the call is made, sharing the data
// across multiple attributes is forbidden.
func (am AttributeMap) InsertBytes(k string, v []byte) {
	if _, existing := am.Get(k); !existing {
		*am.orig = append(*am.orig, newAttributeKeyValueBytes(k, v))
	}
}

// Update updates an existing AnyValue with a value.
// No action is applied to the map where the key does not exist.
//
// Calling this function with a zero-initialized AnyValue struct will cause a panic.
//
// Important: this function should not be used if the caller has access to
// the raw value to avoid an extra allocation.
func (am AttributeMap) Update(k string, v AnyValue) {
	if av, existing := am.Get(k); existing {
		v.copyTo(av.orig)
	}
}

// UpdateString updates an existing string Value with a value.
// No action is applied to the map where the key does not exist.
func (am AttributeMap) UpdateString(k string, v string) {
	if av, existing := am.Get(k); existing {
		av.SetStringVal(v)
	}
}

// UpdateInt updates an existing int Value with a value.
// No action is applied to the map where the key does not exist.
func (am AttributeMap) UpdateInt(k string, v int64) {
	if av, existing := am.Get(k); existing {
		av.SetIntVal(v)
	}
}

// UpdateDouble updates an existing double Value with a value.
// No action is applied to the map where the key does not exist.
func (am AttributeMap) UpdateDouble(k string, v float64) {
	if av, existing := am.Get(k); existing {
		av.SetDoubleVal(v)
	}
}

// UpdateBool updates an existing bool Value with a value.
// No action is applied to the map where the key does not exist.
func (am AttributeMap) UpdateBool(k string, v bool) {
	if av, existing := am.Get(k); existing {
		av.SetBoolVal(v)
	}
}

// UpdateBytes updates an existing []byte Value with a value.
// No action is applied to the map where the key does not exist.
// The caller must ensure the []byte passed in is not modified after the call is made, sharing the data
// across multiple attributes is forbidden.
func (am AttributeMap) UpdateBytes(k string, v []byte) {
	if av, existing := am.Get(k); existing {
		av.SetBytesVal(v)
	}
}

// Upsert performs the Insert or Update action. The AnyValue is
// inserted to the map that did not originally have the key. The key/value is
// updated to the map where the key already existed.
//
// Calling this function with a zero-initialized AnyValue struct will cause a panic.
//
// Important: this function should not be used if the caller has access to
// the raw value to avoid an extra allocation.
func (am AttributeMap) Upsert(k string, v AnyValue) {
	if av, existing := am.Get(k); existing {
		v.copyTo(av.orig)
	} else {
		*am.orig = append(*am.orig, newAttributeKeyValue(k, v))
	}
}

// UpsertString performs the Insert or Update action. The AnyValue is
// inserted to the map that did not originally have the key. The key/value is
// updated to the map where the key already existed.
func (am AttributeMap) UpsertString(k string, v string) {
	if av, existing := am.Get(k); existing {
		av.SetStringVal(v)
	} else {
		*am.orig = append(*am.orig, newAttributeKeyValueString(k, v))
	}
}

// UpsertInt performs the Insert or Update action. The int Value is
// inserted to the map that did not originally have the key. The key/value is
// updated to the map where the key already existed.
func (am AttributeMap) UpsertInt(k string, v int64) {
	if av, existing := am.Get(k); existing {
		av.SetIntVal(v)
	} else {
		*am.orig = append(*am.orig, newAttributeKeyValueInt(k, v))
	}
}

// UpsertDouble performs the Insert or Update action. The double Value is
// inserted to the map that did not originally have the key. The key/value is
// updated to the map where the key already existed.
func (am AttributeMap) UpsertDouble(k string, v float64) {
	if av, existing := am.Get(k); existing {
		av.SetDoubleVal(v)
	} else {
		*am.orig = append(*am.orig, newAttributeKeyValueDouble(k, v))
	}
}

// UpsertBool performs the Insert or Update action. The bool Value is
// inserted to the map that did not originally have the key. The key/value is
// updated to the map where the key already existed.
func (am AttributeMap) UpsertBool(k string, v bool) {
	if av, existing := am.Get(k); existing {
		av.SetBoolVal(v)
	} else {
		*am.orig = append(*am.orig, newAttributeKeyValueBool(k, v))
	}
}

// UpsertBytes performs the Insert or Update action. The []byte Value is
// inserted to the map that did not originally have the key. The key/value is
// updated to the map where the key already existed.
// The caller must ensure the []byte passed in is not modified after the call is made, sharing the data
// across multiple attributes is forbidden.
func (am AttributeMap) UpsertBytes(k string, v []byte) {
	if av, existing := am.Get(k); existing {
		av.SetBytesVal(v)
	} else {
		*am.orig = append(*am.orig, newAttributeKeyValueBytes(k, v))
	}
}

// Sort sorts the entries in the AttributeMap so two instances can be compared.
// Returns the same instance to allow nicer code like:
//   assert.EqualValues(t, expected.Sort(), actual.Sort())
func (am AttributeMap) Sort() AttributeMap {
	// Intention is to move the nil values at the end.
	sort.SliceStable(*am.orig, func(i, j int) bool {
		return (*am.orig)[i].Key < (*am.orig)[j].Key
	})
	return am
}

// Len returns the length of this map.
//
// Because the AttributeMap is represented internally by a slice of pointers, and the data are comping from the wire,
// it is possible that when iterating using "Range" to get access to fewer elements because nil elements are skipped.
func (am AttributeMap) Len() int {
	return len(*am.orig)
}

// Range calls f sequentially for each key and value present in the map. If f returns false, range stops the iteration.
//
// Example:
//
//   sm.Range(func(k string, v AnyValue) bool {
//       ...
//   })
func (am AttributeMap) Range(f func(k string, v AnyValue) bool) {
	for i := range *am.orig {
		kv := &(*am.orig)[i]
		if !f(kv.Key, AnyValue{&kv.Value}) {
			break
		}
	}
}

// CopyTo copies all elements from the current map to the dest.
func (am AttributeMap) CopyTo(dest AttributeMap) {
	newLen := len(*am.orig)
	oldCap := cap(*dest.orig)
	if newLen <= oldCap {
		// New slice fits in existing slice, no need to reallocate.
		*dest.orig = (*dest.orig)[:newLen:oldCap]
		for i := range *am.orig {
			akv := &(*am.orig)[i]
			destAkv := &(*dest.orig)[i]
			destAkv.Key = akv.Key
			AnyValue{&akv.Value}.copyTo(&destAkv.Value)
		}
		return
	}

	// New slice is bigger than exist slice. Allocate new space.
	origs := make([]otlpcommon.KeyValue, len(*am.orig))
	for i := range *am.orig {
		akv := &(*am.orig)[i]
		origs[i].Key = akv.Key
		AnyValue{&akv.Value}.copyTo(&origs[i].Value)
	}
	*dest.orig = origs
}

// AsRaw converts an OTLP AttributeMap to a standard go map
func (am AttributeMap) AsRaw() map[string]interface{} {
	rawMap := make(map[string]interface{})
	am.Range(func(k string, v AnyValue) bool {
		switch v.Type() {
		case AnyValueTypeString:
			rawMap[k] = v.StringVal()
		case AnyValueTypeInt:
			rawMap[k] = v.IntVal()
		case AnyValueTypeDouble:
			rawMap[k] = v.DoubleVal()
		case AnyValueTypeBool:
			rawMap[k] = v.BoolVal()
		case AnyValueTypeBytes:
			rawMap[k] = v.BytesVal()
		case AnyValueTypeEmpty:
			rawMap[k] = nil
		case AnyValueTypeMap:
			rawMap[k] = v.MapVal().AsRaw()
		case AnyValueTypeArray:
			rawMap[k] = v.SliceVal().asRaw()
		}
		return true
	})
	return rawMap
}

// asRaw creates a slice out of a AttributeValueSlice.
func (es AttributeValueSlice) asRaw() []interface{} {
	rawSlice := make([]interface{}, 0, es.Len())
	for i := 0; i < es.Len(); i++ {
		v := es.At(i)
		switch v.Type() {
		case AnyValueTypeString:
			rawSlice = append(rawSlice, v.StringVal())
		case AnyValueTypeInt:
			rawSlice = append(rawSlice, v.IntVal())
		case AnyValueTypeDouble:
			rawSlice = append(rawSlice, v.DoubleVal())
		case AnyValueTypeBool:
			rawSlice = append(rawSlice, v.BoolVal())
		case AnyValueTypeBytes:
			rawSlice = append(rawSlice, v.BytesVal())
		case AnyValueTypeEmpty:
			rawSlice = append(rawSlice, nil)
		default:
			rawSlice = append(rawSlice, "<Invalid array value>")
		}
	}
	return rawSlice
}
