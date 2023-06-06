package dba

import (
	"database/sql"

	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

func (da *DatabaseAccess) InsertCategory(tx *sqlx.Tx, c *Category) (int, error) {

	s := sqrl.Insert("category").Columns("name", "description").
		Values(c.Name, c.Description)

	stmt, params, err := s.ToSql()
	if err != nil {
		return 0, err
	}

	var res sql.Result
	if tx == nil {
		res, err = da.db.Exec(stmt, params...)
	} else {
		res, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
