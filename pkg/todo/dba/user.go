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

func (da *DatabaseAccess) InsertUser(tx *sqlx.Tx, user *User) (int32, error) {
	s := sqrl.Insert("user").Columns("username", "password").Values(user.Username, user.Password)
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertUser")
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
			err1 := errors.WithContextValue(ErrEntityAlreadyExists, "entity", "User")
			return 0, err1
		}
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertUser")
		return 0, err1
	}

	id, err := res.LastInsertId()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "InsertUser")
		return 0, err1
	}

	return int32(id), nil

}

func (da *DatabaseAccess) GetUserById(tx *sqlx.Tx, id int32) (*User, error) {

	s := sqrl.Select("*").From("user").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetUserById")
		return nil, err1
	}

	var u User
	if tx == nil {
		err = da.db.Get(&u, stmt, params...)
	} else {
		err = tx.Get(&u, stmt, params...)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err1 := errors.WithCause(ErrEntityNotFound, err)
			err1 = errors.WithContextValue(err1, "entity", "User")
			err1 = errors.WithContextValue(err1, "entityId", fmt.Sprintf("%d", id))
			return nil, err1
		}
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "GetUserById")
		return nil, err1
	}

	return &u, nil
}

func (da *DatabaseAccess) DeleteUserById(tx *sqlx.Tx, id int32) error {

	s := sqrl.Delete().From("user").Where(sqrl.Eq{"id": id})
	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "DeleteUserById")
		return err1
	}

	if tx == nil {
		_, err = da.db.Exec(stmt, params...)
	} else {
		_, err = tx.Exec(stmt, params...)
	}

	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "DeleteUserById")
		return err1
	}

	return nil
}

func (da *DatabaseAccess) SearchUser(tx *sqlx.Tx, username *string) ([]User, error) {

	s := sqrl.Select("*").From("user")
	if username != nil {
		s = s.Where("username LIKE ?", *username)
	}

	stmt, params, err := s.ToSql()
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "SearchUser")
		return nil, err1
	}

	var us []User
	if tx == nil {
		err = da.db.Select(&us, stmt, params...)
	} else {
		err = tx.Select(&us, stmt, params...)
	}
	if err != nil {
		err1 := errors.WithCause(ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "operation", "SearchUser")
		return nil, err1
	}

	return us, nil
}
