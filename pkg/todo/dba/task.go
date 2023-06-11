package dba

import (
	"database/sql"
	"fmt"
	"todo/pkg/todo/errors"

	"github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
)

func (da *DatabaseAccess) InsertTask(tx *sqlx.Tx, t *Task) (int, error) {
	s := sqrl.Insert("task").Columns("name", "fk_category", "priority", "description").
		Values(t.Name, t.FkCategory, t.Priority, t.Description)

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertTask")
		return 0, err1
	}

	var res sql.Result
	if tx == nil {
		res, err = da.db.Exec(stmt, params...)
	} else {
		res, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertTask")
		return 0, err1
	}

	id, err := res.LastInsertId()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertTask")
		return 0, err1
	}

	return int(id), nil
}

func (da *DatabaseAccess) GetTaskById(tx *sqlx.Tx, id int) (*Task, error) {

	s := sqrl.Select("*").From("task").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetTaskById")
		return nil, err1
	}

	var c Task
	if tx == nil {
		err = da.db.Get(&c, stmt, params...)
	} else {
		err = tx.Get(&c, stmt, params...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err1 := errors.WithCause(ErrEntityNotFound, err)
			err1 = errors.WithContextValue(err1, "entity", "Task")
			err1 = errors.WithContextValue(err1, "entityId", fmt.Sprintf("%d", id))
			return nil, err1
		}
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetTaskById")
		return nil, err1
	}

	return &c, nil
}

func (da *DatabaseAccess) DeleteTaskById(tx *sqlx.Tx, id int) error {

	s := sqrl.Delete().From("task").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "DeleteTaskById")
		return err1
	}

	if tx == nil {
		_, err = da.db.Exec(stmt, params...)
	} else {
		_, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "DeleteTaskById")
		return err1
	}

	return nil
}
