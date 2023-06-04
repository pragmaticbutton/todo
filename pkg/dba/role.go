package dba

import (
	"database/sql"

	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

func (da *DatabaseAccess) InsertRole(tx *sqlx.Tx, role *Role) (int, error) {
	s := sqrl.Insert("role").Columns("name", "description").
		Values(role.Name, role.Description)

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
