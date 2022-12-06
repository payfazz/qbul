// Package qbul provide utility to compose sql string for postgres.
//
// it also can be used for other db software other than postgres, as long it
// also uses $1, $2, $3, ... for their argument placeholder
package qbul

import (
	"reflect"
	"strconv"
)

// A Builder is used to build sql string for postgres.
// The zero value is ready to use.
type Builder struct {
	sql    []byte
	params []any
}

// The sql statement.
func (b *Builder) SQL() string { return string(b.sql) }

// The positional parameters that coresponding with string returned by [SQL] method.
func (b *Builder) Params() []any { return b.params }

// P is shorthand for [Param]
func (b *Builder) P(data any) *Builder { return b.Param(data) }

// Add raw query string to the builder
func (b *Builder) Raw(sql string) *Builder {
	if len(b.sql) != 0 {
		b.sql = append(b.sql, ' ')
	}

	b.sql = append(b.sql, sql...)
	return b
}

// Add argument data to the builder
func (b *Builder) Param(data any) *Builder {
	if len(b.sql) != 0 {
		b.sql = append(b.sql, ' ')
	}

	b.sql = append(b.sql, '$')
	if reflect.TypeOf(data).Comparable() {
		bottom := len(b.params) - 32
		if bottom < 0 {
			bottom = 0
		}
		for i := len(b.params) - 1; i >= bottom; i-- {
			if data == b.params[i] {
				b.sql = append(b.sql, strconv.Itoa(i+1)...)
				return b
			}
		}
	}
	b.params = append(b.params, data)
	b.sql = append(b.sql, strconv.Itoa(len(b.params))...)
	return b
}

// Reset the builder.
func (b *Builder) Reset() *Builder {
	b.sql = b.sql[:0]
	b.params = b.params[:0]
	return b
}
