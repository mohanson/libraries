// Package jany provides a set of functions to parse and interact with json data.
package jany

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strconv"
)

// Data returns a new Jany instance from the given byte slice.
func Data(b []byte) (*Jany, error) {
	return Read(bytes.NewReader(b))
}

// Jany is a struct that holds any type of data, allowing for flexible parsing and manipulation.
type Jany struct {
	j any
}

// Bool returns the bool representation of the current Jany instance.
func (j *Jany) Bool() bool {
	return j.j.(bool)
}

// Dict returns the dict representation of the current Jany instance.
func (j *Jany) Dict() map[string]*Jany {
	a := j.j.(map[string]any)
	r := map[string]*Jany{}
	for k, v := range a {
		r[k] = &Jany{j: v}
	}
	return r
}

// Float32 returns the float32 representation of the current Jany instance.
func (j *Jany) Float32() float32 {
	f, err := strconv.ParseFloat(j.j.(json.Number).String(), 32)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return float32(f)
}

// Float64 returns the float64 representation of the current Jany instance.
func (j *Jany) Float64() float64 {
	f, err := strconv.ParseFloat(j.j.(json.Number).String(), 64)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return float64(f)
}

// Get returns the Jany[k].
func (j *Jany) Get(k string) *Jany {
	a := j.j.(map[string]any)
	return &Jany{j: a[k]}
}

// Idx returns the Jany[i].
func (j *Jany) Idx(k int) *Jany {
	a := j.j.([]any)
	return &Jany{j: a[k]}
}

// Int8 returns the int8 representation of the current Jany instance.
func (j *Jany) Int8() int8 {
	n, err := strconv.ParseInt(j.j.(json.Number).String(), 0, 8)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return int8(n)
}

// Int16 returns the int16 representation of the current Jany instance.
func (j *Jany) Int16() int16 {
	n, err := strconv.ParseInt(j.j.(json.Number).String(), 0, 16)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return int16(n)
}

// Int32 returns the int32 representation of the current Jany instance.
func (j *Jany) Int32() int32 {
	n, err := strconv.ParseInt(j.j.(json.Number).String(), 0, 32)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return int32(n)
}

// Int64 returns the int64 representation of the current Jany instance.
func (j *Jany) Int64() int64 {
	n, err := strconv.ParseInt(j.j.(json.Number).String(), 0, 64)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return int64(n)
}

// Int returns the int representation of the current Jany instance.
func (j *Jany) Int() int {
	n, err := strconv.ParseInt(j.j.(json.Number).String(), 0, 64)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return int(n)
}

// List returns the list representation of the current Jany instance.
func (j *Jany) List() []*Jany {
	a := j.j.([]any)
	r := make([]*Jany, len(a))
	for i, e := range a {
		r[i] = &Jany{j: e}
	}
	return r
}

// String returns the string representation of the current Jany instance.
func (j *Jany) String() string {
	return j.j.(string)
}

// Uint8 returns the uint8 representation of the current Jany instance.
func (j *Jany) Uint8() uint8 {
	n, err := strconv.ParseUint(j.j.(json.Number).String(), 0, 8)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return uint8(n)
}

// Uint16 returns the uint16 representation of the current Jany instance.
func (j *Jany) Uint16() uint16 {
	n, err := strconv.ParseUint(j.j.(json.Number).String(), 0, 16)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return uint16(n)
}

// Uint32 returns the uint32 representation of the current Jany instance.
func (j *Jany) Uint32() uint32 {
	n, err := strconv.ParseUint(j.j.(json.Number).String(), 0, 32)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return uint32(n)
}

// Uint64 returns the uint64 representation of the current Jany instance.
func (j *Jany) Uint64() uint64 {
	n, err := strconv.ParseUint(j.j.(json.Number).String(), 0, 64)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return uint64(n)
}

// Uint returns the uint representation of the current Jany instance.
func (j *Jany) Uint() uint {
	n, err := strconv.ParseUint(j.j.(json.Number).String(), 0, 64)
	if err != nil {
		log.Panicln("jany:", err)
	}
	return uint(n)
}

// Read returns a new Jany instance from the given reader.
func Read(r io.Reader) (*Jany, error) {
	j := new(Jany)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(&j.j)
	return j, err
}
