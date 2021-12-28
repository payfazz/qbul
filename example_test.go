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
	query.Append(
		`select * from people`,
		`where id =`, qbul.Param(id),
		`order by id`,
	)
	if ascOrder {
		query.Append(`asc`)
	} else {
		query.Append(`desc`)
	}

	// do the query
	db.Query(query.SQL(), query.Params()...)
}
