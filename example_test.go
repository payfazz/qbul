package qbul_test

import (
	"database/sql"

	"github.com/payfazz/qbul"
)

func Example() {
	var db *sql.DB // = sql.Open(...)

	// parameter from request
	var ascOrder bool // = ...
	var id int        // = ...

	var query qbul.Builder
	query.Add(
		`select * from people`,
		`where id =`, qbul.P(id),
		`order by id`,
	)
	if ascOrder {
		query.Add(`asc`)
	} else {
		query.Add(`desc`)
	}

	// do the query
	db.Query(query.SQL(), query.Params()...)
}
