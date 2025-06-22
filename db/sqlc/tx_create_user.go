package db

import (
	"context"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreateCallback func(username string) error
}

type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var createUserTxResult CreateUserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		createUserTxResult.User, err = store.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		err = arg.AfterCreateCallback(createUserTxResult.User.Username)
		return err
	})

	return createUserTxResult, err
}
