// Package qbul provide utility to compose sql string for postgres.
//
// it also can be used for other db software other than postgres, as long it
// also uses $1, $2, $3, ... for their argument placeholder
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
	params      []any
	paramsIndex map[any]int
}

// Param is the type returned by P function
type Param struct{ data any }

// P is used to indicate sql parameter.
func P(data any) Param { return Param{data} }

// SQL return the complete sql statement.
func (b *Builder) SQL() string { return b.sql.String() }

// Params return positional parameters that coresponding with string returned by SQL method.
func (b *Builder) Params() []any { return b.params }

// Add data into builder.
// data must be string or Param returned by P function.
//
// CAUTION: you must use P function to pass string parameter, without it, the string is appended
// to query directly, if you do that, you are vulnerable to sql injection.
func (b *Builder) Add(data ...any) *Builder {
	for _, item := range data {
		if b.sql.Len() != 0 {
			b.sql.WriteByte(' ')
		}

		switch item := item.(type) {
		case string:
			b.sql.WriteString(item)

		case Param:
			data := item.data
			pos := len(b.params) + 1

			if reflect.TypeOf(data).Comparable() {
				if cachedPos, ok := b.paramsIndex[data]; ok {
					pos = cachedPos
				} else {
					if b.paramsIndex == nil {
						b.paramsIndex = make(map[any]int)
					}
					b.paramsIndex[data] = pos
					b.params = append(b.params, data)
				}
			} else {
				b.params = append(b.params, data)
			}

			b.sql.WriteByte('$')
			b.sql.WriteString(strconv.Itoa(pos))
		default:
			panic(`invalid argument: must be a value with type "string" or "Param"`)
		}
	}
	return b
}

// Reset the builder.
//
// will pass the arguments to Add method.
func (b *Builder) Reset(data ...any) *Builder {
	b.sql.Reset()
	b.params = nil
	b.paramsIndex = nil
	return b.Add(data...)
}
