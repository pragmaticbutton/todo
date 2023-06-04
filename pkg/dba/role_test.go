package dba

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func prepareRole() int {

	var id int
	err := da.ExecuteInTransaction(func(tx *sqlx.Tx) error {
		role := Role{Name: "role_name", Description: sql.NullString{Valid: true, String: "role description"}}
		var err1 error
		id, err1 = da.InsertRole(tx, &role)
		if err1 != nil {
			return err1
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return id
}
