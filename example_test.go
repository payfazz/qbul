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
	query.Reset().
		Raw(`select * from people`).
		Raw(`where id =`).P(id).
		Raw(`order by id`)

	if ascOrder {
		query.Raw(`asc`)
	} else {
		query.Raw(`desc`)
	}

	// do the query
	db.Query(query.SQL(), query.Params()...)
}
