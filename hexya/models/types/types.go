// Copyright 2016 NDP Systèmes. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"database/sql/driver"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/hexya-erp/hexya/hexya/tools/logging"
)

var log *logging.Logger

// A Context is a map of objects that is passed along from function to function
// during a transaction. A Context is read only.
type Context struct {
	values map[string]interface{}
}

// Copy returns a shallow copy of the Context
func (c Context) Copy() *Context {
	newCtx := NewContext()
	for k, v := range c.values {
		newCtx.values[k] = v
	}
	return newCtx
}

// Get returns the value of the given key in this Context
func (c *Context) Get(key string) interface{} {
	value := c.values[key]
	return value
}

// GetString returns the value of the given key in
// this Context as a string.
// It panics if the value is not of type string
func (c *Context) GetString(key string) string {
	return c.Get(key).(string)
}

// GetInteger returns the value of the given key in
// this Context as an int64.
// It panics if the value cannot be casted to int64
func (c *Context) GetInteger(key string) int64 {
	val := c.Get(key)
	res := castToInteger(val, key)
	return res
}

// GetFloat returns the value of the given key in
// this Context as a float64.
// It panics if the value cannot be casted to float64
func (c *Context) GetFloat(key string) float64 {
	val := c.Get(key)
	res := castToFloat(val, key)
	return res
}

// GetStringSlice returns the value of the given key in
// this Context as a []string.
// It panics if the value is not a slice or if any value
// is not a string
func (c *Context) GetStringSlice(key string) []string {
	val := c.Get(key)
	var res []string
	switch value := val.(type) {
	case []string:
		res = value
	case []interface{}:
		res = make([]string, len(value))
		for i, v := range value {
			res[i] = v.(string)
		}
	}
	return res
}

// GetIntegerSlice returns the value of the given key in
// this Context as a []int64.
// It panics if the value is not a slice or if any value
// cannot be casted to int64
func (c *Context) GetIntegerSlice(key string) []int64 {
	val := c.Get(key)
	rVal := reflect.ValueOf(val)
	if rVal.Kind() != reflect.Slice {
		log.Panic("Value in Context is not a slice", "key", key, "value", val)
	}
	res := make([]int64, rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		res[i] = castToInteger(rVal.Index(i).Interface(), key)
	}
	return res
}

// GetFloatSlice returns the value of the given key in
// this Context as a []float64.
// It panics if the value is not a slice or if any value
// cannot be casted to float64
func (c *Context) GetFloatSlice(key string) []float64 {
	val := c.Get(key)
	rVal := reflect.ValueOf(val)
	if rVal.Kind() != reflect.Slice {
		log.Panic("Value in Context is not a slice", "key", key, "value", val)
	}
	res := make([]float64, rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		res[i] = castToFloat(rVal.Index(i).Interface(), key)
	}
	return res
}

// castToInteger casts the given val to int64 if it is
// a number type. Panics otherwise
func castToInteger(val interface{}, key string) int64 {
	var res int64
	switch value := val.(type) {
	case int64:
		res = value
	case int, int8, int16, int32, uint, uint8, uint16, uint32, uint64, float32, float64:
		res, _ = strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
	default:
		log.Panic("Context value cannot be cast to int64", "key", key, "value", val)
	}
	return res
}

// castToFloat casts the given val to float64 if it is
// a number type. Panics otherwise
func castToFloat(val interface{}, key string) float64 {
	var res float64
	switch value := val.(type) {
	case float64:
		res = value
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32:
		res, _ = strconv.ParseFloat(fmt.Sprintf("%d", value), 64)
	default:
		log.Panic("Context value cannot be cast to float64", "key", key, "value", val)
	}
	return res
}

// HasKey returns true if this Context has the given key
func (c *Context) HasKey(key string) bool {
	_, exists := c.values[key]
	return exists
}

// WithKey returns a copy of this context with the given key/value.
// If key already exists, it is overwritten.
func (c Context) WithKey(key string, value interface{}) *Context {
	c.values[key] = value
	return &c
}

// IsEmpty returns true if this Context has no entries.
func (c Context) IsEmpty() bool {
	if len(c.values) == 0 {
		return true
	}
	return false
}

// ToMap returns a copy of the map of values of this context
func (c Context) ToMap() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range c.values {
		res[k] = v
	}
	return res
}

// UnmarshalXMLAttr is the XML unmarshalling method of Context.
func (c *Context) UnmarshalXMLAttr(attr xml.Attr) error {
	var cm map[string]interface{}
	err := json.Unmarshal([]byte(attr.Value), &cm)
	(*c).values = cm
	return err
}

// MarshalJSON method for Context
func (c *Context) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.values)
}

// UnmarshalJSON method for Context
func (c *Context) UnmarshalJSON(data []byte) error {
	var cm map[string]interface{}
	err := json.Unmarshal(data, &cm)
	(*c).values = cm
	return err
}

var _ xml.UnmarshalerAttr = &Context{}
var _ json.Marshaler = &Context{}
var _ json.Unmarshaler = &Context{}

// NewContext returns a new Context instance
func NewContext(data ...map[string]interface{}) *Context {
	var values map[string]interface{}
	if len(data) > 0 {
		values = data[0]
	} else {
		values = make(map[string]interface{})
	}
	return &Context{
		values: values,
	}
}

// Digits holds precision and scale information for a float (numeric) type:
// - The precision: the total number of digits
// - The scale: the number of digits to the right of the decimal point
// (PostgresSQL definitions)
type Digits struct {
	Precision int8
	Scale     int8
}

// Date type that JSON marshal and unmarshals as "YYYY-MM-DD"
type Date time.Time

// IsNull returns true if the Date is the zero value
func (d Date) IsNull() bool {
	if time.Time(d).Format("2006-01-02") == "0001-01-01" {
		return true
	}
	return false
}

// MarshalJSON for Date type
func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsNull() {
		return []byte("false"), nil
	}
	dateStr := time.Time(d).Format("2006-01-02")
	dateStr = fmt.Sprintf(`"%s"`, dateStr)
	return []byte(dateStr), nil
}

// Value formats our Date for storing in database
// Especially handles empty Date.
func (d Date) Value() (driver.Value, error) {
	if d.IsNull() {
		return driver.Value("0001-01-01"), nil
	}
	return driver.Value(d), nil
}

// Today returns the current date
func Today() Date {
	return Date(time.Now())
}

// DateTime type that JSON marshals and unmarshals as "YYYY-MM-DD HH:MM:SS"
type DateTime time.Time

// IsNull returns true if the DateTime is the zero value
func (d DateTime) IsNull() bool {
	if time.Time(d).Format("2006-01-02 15:04:05") == "0001-01-01 00:00:00" {
		return true
	}
	return false
}

// Now returns the current date/time
func Now() DateTime {
	return DateTime(time.Now())
}

// MarshalJSON for DateTime type
func (d DateTime) MarshalJSON() ([]byte, error) {
	if d.IsNull() {
		return []byte("false"), nil
	}
	dateStr := time.Time(d).Format("2006-01-02 15:04:05")
	dateStr = fmt.Sprintf(`"%s"`, dateStr)
	return []byte(dateStr), nil
}

// Value formats our DateTime for storing in database
// Especially handles empty DateTime.
func (d DateTime) Value() (driver.Value, error) {
	if d.IsNull() {
		return driver.Value("0001-01-01 00:00:00"), nil
	}
	return driver.Value(time.Time(d).Format("2006-01-02 15:04:05")), nil
}

// A Selection is a set of possible (key, label) values for a model
// "selection" field.
type Selection map[string]string

// MarshalJSON function for the Selection type
func (s Selection) MarshalJSON() ([]byte, error) {
	var selSlice [][2]string
	for key, value := range s {
		selSlice = append(selSlice, [2]string{0: key, 1: value})
	}
	return json.Marshal(selSlice)
}

var _ json.Marshaler = Selection{}

func init() {
	log = logging.GetLogger("types")
}