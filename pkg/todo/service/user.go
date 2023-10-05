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

func (svc *toDoService) GetUser(ctx context.Context, id int32) (*restapi.UserOut, error) {

	u, err := svc.da.GetUserById(nil, id)
	if err != nil {
		return nil, err
	}

	out := dbaToRestUserOut(u)
	return out, nil
}

func (svc *toDoService) DeleteUser(ctx context.Context, id int32) error {

	err := svc.da.ExecuteInTransaction(func(tx *sqlx.Tx) error {

		_, err := svc.da.GetUserById(tx, id)
		if err != nil {
			return err
		}

		err = svc.da.DeleteUserById(tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (svc *toDoService) SearchUser(ctx context.Context, params *restapi.SearchUserParams) (*restapi.SearchUserOut, error) {

	var username *string
	if params.Username != nil {
		n := replaceWildCards(*params.Username)
		username = &n
	}

	us, err := svc.da.SearchUser(nil, username)
	if err != nil {
		return nil, err
	}

	usOut := make([]restapi.UserOut, len(us))
	for i, u := range us {
		usOut[i] = *dbaToRestUserOut(&u)
	}

	out := &restapi.SearchUserOut{Users: &usOut}

	return out, nil
}
