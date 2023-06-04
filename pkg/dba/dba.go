package dba

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
		return err
	}

	err = f(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("rollback failed")
		}
		return err
	}

	tx.Commit()

	return nil
}
