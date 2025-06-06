// Copyright 2025 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package playlist

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalidQuotedString   = errors.New("invalid quoted string")
	errInvalidUnquotedString = errors.New("invalid unquoted string")
	errInvalidDecimalHexSeq  = errors.New("invalid hexadecimal sequence")
	errInvalidDecimalInteger = errors.New("invalid decimal integer")
	errInvalidDecimalFloat   = errors.New("invalid decimal float")
	errInvalidBool           = errors.New("invalid bool")
	errInvalidTime           = errors.New("invalid time")

	errInvalidAttribute      = errors.New("invalid attribute")
	errInvalidAttributeName  = errors.New("invalid attribute name")
	errInvalidAttributeValue = errors.New("invalid attribute value")
)

type _Encoder interface {
	encode(w io.Writer) (err error)
}

type _QuotedString string

func (v _QuotedString) IsZero() bool { return v == "" }

func (v _QuotedString) valid() bool {
	return !strings.ContainsAny(string(v), "\r\n\"")
}

func (v _QuotedString) get() string { return string(v) }

func (v _QuotedString) encode(w io.Writer) error {
	if !v.valid() {
		return errInvalidQuotedString
	}

	_, err := io.WriteString(w, strconv.Quote(string(v)))
	return err
}

func (v *_QuotedString) decode(s string) error {
	value, err := strconv.Unquote(s)
	if err != nil || value == "" {
		return errInvalidQuotedString
	}

	if *v = _QuotedString(value); !v.valid() {
		return errInvalidQuotedString
	}

	return nil
}

type _UnquotedString string

func (v _UnquotedString) IsZero() bool { return v == "" }

func (v _UnquotedString) valid() bool {
	return v != "" && !strings.ContainsAny(string(v), `, "`) // Comma/Space/Quote
}

func (v _UnquotedString) get() string { return string(v) }

func (v _UnquotedString) encode(w io.Writer) error {
	if !v.valid() {
		return errInvalidUnquotedString
	}

	_, err := io.WriteString(w, string(v))
	return err
}

func (v *_UnquotedString) decode(s string) error {
	if s == "" {
		return errInvalidUnquotedString
	}

	if *v = _UnquotedString(s); !v.valid() {
		return errInvalidUnquotedString
	}
	return nil
}

type _DecimalInteger uint64

func (v _DecimalInteger) IsZero() bool { return v == 0 }

func (v _DecimalInteger) get() uint64 { return uint64(v) }

func (v _DecimalInteger) encode(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%d", uint64(v))
	return err
}

func (v *_DecimalInteger) decode(s string, min uint64) (err error) {
	value, err := strconv.ParseUint(s, 10, 64)
	if err != nil || value < min {
		return errInvalidDecimalInteger
	}
	*v = _DecimalInteger(value)
	return
}

type _DecimalFloat float64

func (v _DecimalFloat) IsZero() bool { return v == 0 }
func (v _DecimalFloat) valid() bool  { return v >= 0 }

func (v _DecimalFloat) get() float64 { return float64(v) }

func (v _DecimalFloat) encode(w io.Writer) error {
	if !v.valid() {
		return errInvalidDecimalFloat
	}

	s := strconv.FormatFloat(float64(v), 'f', -1, 64)
	if index := strings.IndexByte(s, '.'); index >= 0 && len(s[index+1:]) > 3 {
		s = s[:index+4] // Limit to 3 decimal places
	}

	_, err := io.WriteString(w, s)
	return err
}

func (v *_DecimalFloat) decode(s string) (err error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil || !v.valid() {
		return errInvalidDecimalFloat
	}

	*v = _DecimalFloat(value)
	return
}

type _SignDecimalFloat float64

func (v _SignDecimalFloat) IsZero() bool { return v == 0 }

func (v _SignDecimalFloat) get() float64 { return float64(v) }

func (v _SignDecimalFloat) encode(w io.Writer) error {
	_, err := io.WriteString(w, strconv.FormatFloat(float64(v), 'f', -1, 64))
	return err
}

func (v *_SignDecimalFloat) decode(s string) (err error) {
	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errInvalidDecimalFloat
	}

	*v = _SignDecimalFloat(value)
	return
}

type _HexSequence []byte

func (v _HexSequence) IsZero() bool { return len(v) == 0 }

func (v _HexSequence) valid() bool { return len(v) > 0 }

func (v _HexSequence) String() string {
	if !v.valid() {
		return ""
	}

	var b strings.Builder
	b.Grow(hex.EncodedLen(len(v)) + 2)
	_, _ = hex.NewEncoder(&b).Write(v)
	return b.String()
}

func (v _HexSequence) encode(w io.Writer) error {
	if len(v) == 0 {
		return errInvalidDecimalHexSeq
	}

	_, err := fmt.Fprintf(w, "0x%02X", []byte(v))
	return err
}

func (v *_HexSequence) decode(s string) (err error) {
	if len(s) < 3 || (s[:2] != "0x" && s[:2] != "0X") {
		return errInvalidDecimalHexSeq
	}

	value, err := hex.DecodeString(s[2:])
	if err != nil {
		return errInvalidDecimalHexSeq
	}

	*v = _HexSequence(value)
	return
}

type _Bool bool

func (v _Bool) IsZero() bool { return !bool(v) }

func (v _Bool) get() bool { return bool(v) }

func (v _Bool) encode(w io.Writer) error {
	var s string
	if v {
		s = "YES"
	} else {
		s = "NO"
	}
	_, err := io.WriteString(w, s)
	return err
}

func (v *_Bool) decode(s string) (err error) {
	switch s {
	case "YES":
		*v = true
	case "NO":
		*v = false
	default:
		err = errInvalidBool
	}
	return
}

const _TimeLayout = "2006-01-02T15:04:05.999Z07:00"

type _Time time.Time

func (v _Time) IsZero() bool { return time.Time(v).IsZero() }

func (v _Time) get() time.Time { return time.Time(v) }

func (v _Time) encode(w io.Writer) (err error) {
	if time.Time(v).IsZero() {
		return errInvalidTime
	}

	s := time.Time(v).Format(_TimeLayout)
	_, err = io.WriteString(w, s)
	return
}

func (v *_Time) decode(s string) (err error) {
	t, err := time.Parse(_TimeLayout, s)
	if err != nil {
		return fmt.Errorf("%w: %w", errInvalidTime, err)
	} else if t.IsZero() {
		return errInvalidTime
	} else {
		*v = _Time(t)
	}
	return
}

type _EnumValuer interface {
	~string
	validate() error
}

type _Enum[T _EnumValuer] struct {
	value T
}

func newEnum[T _EnumValuer](value T) _Enum[T] {
	return _Enum[T]{value: value}
}

func (v _Enum[T]) IsZero() bool { return v.value == "" }

func (v _Enum[T]) get() T { return v.value }

func (v _Enum[T]) encode(w io.Writer) error {
	if err := v.value.validate(); err != nil {
		return err
	}

	_, err := io.WriteString(w, string(v.value))
	return err
}

func (v *_Enum[T]) decode(s string) error {
	v.value = T(s)
	return v.value.validate()
}

type _Value interface {
	_Encoder
	IsZero() bool
}

type _Attr struct {
	Name  _UnquotedString
	Value _Value
}

func _NewAttr(name string, value _Value) _Attr {
	return _Attr{Name: _UnquotedString(name), Value: value}
}

func (v _Attr) IsZero() bool { return v.Value.IsZero() }

func (v _Attr) encode(w io.Writer) (err error) {
	if err = v.Name.encode(w); err != nil {
		return
	}

	if _, err = io.WriteString(w, "="); err != nil {
		return
	}

	err = v.Value.encode(w)
	return
}
