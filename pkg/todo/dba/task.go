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

func (da *DatabaseAccess) InsertTask(tx *sqlx.Tx, t *Task) (int32, error) {
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
		var me *mysql.MySQLError
		if ok := stderrors.As(err, &me); ok && me.Number == 1062 {
			err1 := errors.WithContextValue(ErrEntityAlreadyExists, "entity", "Task")
			return 0, err1
		}
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

	return int32(id), nil
}

func (da *DatabaseAccess) GetTaskById(tx *sqlx.Tx, id int32) (*Task, error) {

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

func (da *DatabaseAccess) DeleteTaskById(tx *sqlx.Tx, id int32) error {

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

func (da *DatabaseAccess) SearchTask(tx *sqlx.Tx, name *string, fkCategory *int32, priority *TaskPriorityType, done *int8, p *Pagination) ([]Task, error) {

	s := sqrl.Select("*").From("task")
	if name != nil {
		s = s.Where("name LIKE ?", *name)
	}

	if fkCategory != nil {
		s = s.Where(sqrl.Eq{"fk_category": *fkCategory})
	}

	if priority != nil {
		s = s.Where(sqrl.Eq{"priority": *priority})
	}

	if done != nil {
		s = s.Where(sqrl.Eq{"done": *done})
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
		err1 = errors.WithContextValue(err1, "operation", "SearchTask")
		return nil, err1
	}

	var ts []Task
	if tx == nil {
		err = da.db.Select(&ts, stmt, params...)
	} else {
		err = tx.Select(&ts, stmt, params...)
	}
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "SearchTask")
		return nil, err1
	}

	return ts, nil
}

func (da *DatabaseAccess) CountTask(tx *sqlx.Tx, name *string, fkCategory *int32, priority *TaskPriorityType, done *int8) (int32, error) {

	s := sqrl.Select("COUNT(*)").From("task")
	if name != nil {
		s = s.Where("name LIKE ?", *name)
	}

	if fkCategory != nil {
		s = s.Where(sqrl.Eq{"fk_category": *fkCategory})
	}

	if priority != nil {
		s = s.Where(sqrl.Eq{"priority": *priority})
	}

	if done != nil {
		s = s.Where(sqrl.Eq{"done": *done})
	}

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "CountTask")
		return 0, err1
	}

	var c int
	if tx == nil {
		err = da.db.Get(&c, stmt, params...)
	} else {
		err = tx.Get(&c, stmt, params...)
	}
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "CountTask")
		return 0, err1
	}

	return int32(c), nil
}

func (da *DatabaseAccess) UpdateTask(tx *sqlx.Tx, t *Task) error {

	s := sqrl.Update("task").Set("name", t.Name).Set("description", t.Description).
		Set("priority", t.Priority).Set("done", t.Done).Set("fk_category", t.FkCategory).
		Where(sqrl.Eq{"id": t.Id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "UpdateTask")
		return err1
	}

	if tx == nil {
		_, err = da.db.Exec(stmt, params...)
	} else {
		_, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "UpdateTask")
		return err1
	}

	return nil
}
