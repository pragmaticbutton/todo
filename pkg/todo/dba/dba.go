package dba

import (
	"net/http"
	"todo/pkg/todo/errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	ErrDatabaseError  = errors.ToDoError{ErrorCode: errors.ERROR_CODE_DATABASE_ERROR, Text: "Database operation failed", HttpStatus: http.StatusInternalServerError}
	ErrEntityNotFound = errors.ToDoError{ErrorCode: errors.ERROR_CODE_ENTITY_NOT_FOUND, Text: "Entity not found", HttpStatus: http.StatusNotFound}
)

type DatabaseAccess struct {
	db *sqlx.DB
}

func NewDatabaseAccess(dsn string) (*DatabaseAccess, error) {
	db, err := openDB(dsn)
	if err != nil {
		return nil, err
	}

	da := DatabaseAccess{db: db}
	return &da, nil
}

func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (da *DatabaseAccess) ExecuteInTransaction(f func(tx *sqlx.Tx) error) error {
	tx, err := da.db.Beginx()
	if err != nil {
		return errors.WithCause(&ErrDatabaseError, err)
	}

	err = f(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			err1 := errors.WithCause(&ErrDatabaseError, err)
			err1 = errors.WithContextValue(err1, "reason", "rollback failed")
			return err1
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		err1 := errors.WithCause(&ErrDatabaseError, err)
		err1 = errors.WithContextValue(err1, "reason", "commit failed")
		return err1
	}

	return nil
}
