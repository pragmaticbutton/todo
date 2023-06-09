package dba

import (
	"database/sql"
	"fmt"
	"todo/pkg/todo/errors"

	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

func (da *DatabaseAccess) InsertCategory(tx *sqlx.Tx, c *Category) (int, error) {

	s := sqrl.Insert("category").Columns("name", "description").
		Values(c.Name, c.Description)

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertCategory")
		return 0, err1
	}

	var res sql.Result
	if tx == nil {
		res, err = da.db.Exec(stmt, params...)
	} else {
		res, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertCategory")
		return 0, err1
	}

	id, err := res.LastInsertId()
	if err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertCategory")
		return 0, err1
	}

	return int(id), nil
}

func (da *DatabaseAccess) GetCategoryById(tx *sqlx.Tx, id int) (*Category, error) {

	s := sqrl.Select("*").From("category").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetCategoryById")
		return nil, err1
	}

	var c Category
	if tx == nil {
		err = da.db.Get(&c, stmt, params...)
	} else {
		err = tx.Get(&c, stmt, params...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err1 := errors.WithCause(&ErrEntityNotFound, err)
			err1 = errors.WithContextValue(err1, "entity", "Category")
			err1 = errors.WithContextValue(err1, "entityId", fmt.Sprintf("%d", id))
			return nil, err1
		}
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetCategoryById")
		return nil, err1
	}

	return &c, nil
}

func (da *DatabaseAccess) SearchCategory(tx *sqlx.Tx, name *string) ([]Category, error) {

	s := sqrl.Select("*").From("category")
	if name != nil {
		s = s.Where("name LIKE ?", name)
	}

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "SearchCategory")
		return nil, err1
	}

	var cs []Category
	if tx == nil {
		err = da.db.Select(&cs, stmt, params...)
	} else {
		err = tx.Select(&cs, stmt, params...)
	}
	if err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "SearchCategory")
		return nil, err1
	}

	return cs, nil
}
