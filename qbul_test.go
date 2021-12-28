package qbul

import (
	"reflect"
	"testing"
	"time"
)

func TestNormalUsage(t *testing.T) {
	t.Parallel()

	q := Build(
		`select * from people`,
		`where id =`, Param(10),
		`and name like`, Param("Bob%"),
		`order by id asc`,
	)

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
}

func TestReuseParam(t *testing.T) {
	t.Parallel()

	now := time.Now()
	q := Build(
		`select * from people`,
		`where birth_time <=`, Param(now), `and`, Param(now), `<= death_time`,
		`and name like`, Param("Bob%"),
	)

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

	q := Build(
		`select * from people`,
		`where id in (select * from unnest(`, Param([]int{1, 2, 3}), `::int4[]))`,
	)

	if q.SQL() != `select * from people where id in (select * from unnest( $1 ::int4[]))` {
		t.Fatalf("invalid sql")
	}

	p := q.Params()
	if len(p) != 1 {
		t.Fatalf("invalid params length")
	}

	if !reflect.DeepEqual(p[0], []int{1, 2, 3}) {
		t.Fatalf("invalid params")
	}
}

func TestInvalidParam(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatalf("should panic")
		}
	}()

	Build(
		`select * from people`,
		`where id =`, 10,
		`and name like`, "Bob%",
		`order by id asc`,
	)
}
