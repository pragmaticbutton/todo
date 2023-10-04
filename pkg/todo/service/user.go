package service

import (
	"context"
	"todo/pkg/todo/restapi"

	"github.com/jmoiron/sqlx"
)

func (svc *toDoService) CreateUser(ctx context.Context, in *restapi.CreateUserIn) (*restapi.UserOut, error) {

	u := restToDbaCreateUserIn(in)

	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		id, err := svc.da.InsertUser(tx, u)
		if err != nil {
			return err
		}

		u, err = svc.da.GetUserById(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	out := dbaToRestUserOut(u)

	return out, nil
}
