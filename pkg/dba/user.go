package dba

import (
	"database/sql"

	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

// InsertUser inserts new record into user table.
func (da *DatabaseAccess) InsertUser(tx *sqlx.Tx, u *User) (int, error) {
	s := sqrl.Insert("user").Columns("username", "fk_role").
		Values(u.Username, u.FkRole)

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
