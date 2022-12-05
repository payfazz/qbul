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

// SQL return the complete sql statement.
func (b *Builder) SQL() string { return b.sql.String() }

// Params return positional parameters that coresponding with string returned by SQL method.
func (b *Builder) Params() []any { return b.params }

// P is shorthand for [Param]
func (b *Builder) P(data any) *Builder { return b.Param(data) }

// Add raw query string to the builder
func (b *Builder) Raw(sql string) *Builder {
	if b.sql.Len() != 0 {
		b.sql.WriteByte(' ')
	}

	b.sql.WriteString(sql)
	return b
}

// Add argument data to the builder
func (b *Builder) Param(data any) *Builder {
	if b.sql.Len() != 0 {
		b.sql.WriteByte(' ')
	}

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

	return b
}

// Reset the builder.
//
// will pass the arguments to Add method.
func (b *Builder) Reset() *Builder {
	b.sql.Reset()
	b.params = nil
	b.paramsIndex = nil
	return b
}
