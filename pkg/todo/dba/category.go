package dba

import (
	"database/sql"
	stderrors "errors"
	"fmt"
	"todo/pkg/todo/errors"

	"github.com/elgris/sqrl"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func (da *DatabaseAccess) InsertCategory(tx *sqlx.Tx, c *Category) (int32, error) {

	s := sqrl.Insert("category").Columns("name", "description").
		Values(c.Name, c.Description)

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
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
		var me *mysql.MySQLError
		if ok := stderrors.As(err, &me); ok && me.Number == 1062 {
			err1 := errors.WithContextValue(ErrEntityAlreadyExists, "entity", "Category")
			return 0, err1
		}
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertCategory")
		return 0, err1
	}

	id, err := res.LastInsertId()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertCategory")
		return 0, err1
	}

	return int32(id), nil
}

func (da *DatabaseAccess) GetCategoryById(tx *sqlx.Tx, id int32) (*Category, error) {

	s := sqrl.Select("*").From("category").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
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
			err1 := errors.WithCause(ErrEntityNotFound, err)
			err1 = errors.WithContextValue(err1, "entity", "Category")
			err1 = errors.WithContextValue(err1, "entityId", fmt.Sprintf("%d", id))
			return nil, err1
		}
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetCategoryById")
		return nil, err1
	}

	return &c, nil
}

func (da *DatabaseAccess) SearchCategory(tx *sqlx.Tx, name *string, p *Pagination) ([]Category, error) {

	s := sqrl.Select("*").From("category")
	if name != nil {
		s = s.Where("name LIKE ?", *name)
	}

	if p != nil {
		var od string
		if p.orderDirection != nil {
			od = *p.orderDirection
		} else {
			od = "asc"
		}
		if p.orderBy != nil {
			s = s.OrderBy(*p.orderBy + " " + od)
		} else {
			s = s.OrderBy("id" + " " + od)
		}
		if p.startIndex != nil {
			s = s.Offset(uint64(*p.startIndex))
		}
		if p.recordsPerPage != nil {
			s = s.Limit(uint64(*p.recordsPerPage))
		} else {
			// mysql doesn't support offset without limit, so this is workaround around that
			s = s.Limit(uint64(18446744073709551615))
		}
	}

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
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
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "SearchCategory")
		return nil, err1
	}

	return cs, nil
}

func (da *DatabaseAccess) CountCategory(tx *sqlx.Tx, name *string) (int32, error) {

	s := sqrl.Select("COUNT(*)").From("category")
	if name != nil {
		s = s.Where("name LIKE ?", *name)
	}

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "CountCategory")
		return 0, err1
	}

	var res int
	if tx == nil {
		err = da.db.Get(&res, stmt, params...)
	} else {
		err = tx.Get(&res, stmt, params...)
	}
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "CountCategory")
		return 0, err1
	}

	return int32(res), nil
}

func (da *DatabaseAccess) DeleteCategoryById(tx *sqlx.Tx, id int32) error {

	s := sqrl.Delete().From("category").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "DeleteCategoryById")
		return err1
	}

	if tx == nil {
		_, err = da.db.Exec(stmt, params...)
	} else {
		_, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "DeleteCategoryById")
		return err1
	}

	return nil
}

func (da *DatabaseAccess) UpdateCategory(tx *sqlx.Tx, c *Category) error {

	s := sqrl.Update("category").Set("name", c.Name).Set("description", c.Description).Where(sqrl.Eq{"id": c.Id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "UpdateCategory")
		return err1
	}

	if tx == nil {
		_, err = da.db.Exec(stmt, params...)
	} else {
		_, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "UpdateCategory")
		return err1
	}

	return nil
}
