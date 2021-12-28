// Package qbul provide utility to compose sql string for postgres.
package qbul

import (
	"reflect"
	"strconv"
	"strings"
)

// A Builder is used to build sql string for postgres.
// The zero value is ready to use.
type Builder struct {
	sql         strings.Builder
	params      []interface{}
	paramsIndex map[interface{}]int
}

// Build will create new Builder and passing the data to its Append method.
func Build(data ...interface{}) *Builder {
	var b Builder
	b.Append(data...)
	return &b
}

type param struct{ data interface{} }

// Param is used to indicate sql parameter.
func Param(data interface{}) interface{} { return param{data} }

// SQL return the complete sql statement.
func (b *Builder) SQL() string { return b.sql.String() }

// Params return positional parameters that coresponding with string returned by SQL method.
func (b *Builder) Params() []interface{} { return b.params }

// Append data into builder.
// data must be string or interface{} returned by Param function.
//
// NOTE: you must use Param function to pass string parameter, without it, the string is appended
// to query directly, if you do that, you are vulnerable to sql injection.
func (b *Builder) Append(data ...interface{}) *Builder {
	for _, item := range data {
		switch x := item.(type) {
		case string:
			if b.sql.Len() != 0 {
				b.sql.WriteByte(' ')
			}
			b.sql.WriteString(x)

		case param:
			p := x.data
			pos := len(b.params) + 1

			if reflect.TypeOf(p).Comparable() {
				if b.paramsIndex == nil {
					b.paramsIndex = make(map[interface{}]int)
				}

				if cachedPos, ok := b.paramsIndex[p]; ok {
					pos = cachedPos
				} else {
					b.paramsIndex[p] = pos
					b.params = append(b.params, p)
				}
			} else {
				b.params = append(b.params, p)
			}

			if b.sql.Len() != 0 {
				b.sql.WriteByte(' ')
			}
			b.sql.WriteByte('$')
			b.sql.WriteString(strconv.Itoa(pos))
		default:
			panic("invalid argument: can't process value with type: " + reflect.TypeOf(item).String())
		}
	}
	return b
}
