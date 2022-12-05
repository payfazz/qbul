package qbul_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/payfazz/qbul"
)

func TestNormalUsage(t *testing.T) {
	t.Parallel()
	var q qbul.Builder
	q.
		Raw(`select * from people`).
		Raw(`where id =`).P(10).
		Raw(`and name like`).P("Bob%").
		Raw(`order by id asc`)

	if q.SQL() != `select * from people where id = $1 and name like $2 order by id asc` {
		t.Fatalf("invalid sql")
	}

	p := q.Params()
	if len(p) != 2 {
		t.Fatalf("invalid params length")
	}

	if p[0] != 10 || p[1] != "Bob%" {
		t.Fatalf("invalid params")
	}

	q.Reset()
	if q.SQL() != "" || len(q.Params()) != 0 {
		t.Fatalf("invalid Reset")
	}
}

func TestReuseParam(t *testing.T) {
	t.Parallel()

	now := time.Now()
	var q qbul.Builder
	q.
		Raw(`select * from people`).
		Raw(`where birth_time <=`).P(now).Raw(`and`).P(now).Raw(`<= death_time`).
		Raw(`and name like`).P("Bob%")

	if q.SQL() != `select * from people where birth_time <= $1 and $1 <= death_time and name like $2` {
		t.Fatalf("invalid sql")
	}

	p := q.Params()
	if len(p) != 2 {
		t.Fatalf("invalid params length")
	}

	if p[0] != now || p[1] != "Bob%" {
		t.Fatalf("invalid params")
	}
}

func TestNonComparableParam(t *testing.T) {
	t.Parallel()

	var q qbul.Builder
	q.
		Raw(`select * from people`).
		Raw(`where id = any(`).P([]int{1, 2, 3}).Raw(`::int4[]) and id = any(`).P([]int{1, 2, 3}).Raw(`::int4[])`)

	if q.SQL() != `select * from people where id = any( $1 ::int4[]) and id = any( $2 ::int4[])` {
		t.Fatalf("invalid sql")
	}

	p := q.Params()
	if len(p) != 2 {
		t.Fatalf("invalid params length")
	}

	if !reflect.DeepEqual(p[0], []int{1, 2, 3}) {
		t.Fatalf("invalid params")
	}
	if !reflect.DeepEqual(p[1], []int{1, 2, 3}) {
		t.Fatalf("invalid params")
	}
}
