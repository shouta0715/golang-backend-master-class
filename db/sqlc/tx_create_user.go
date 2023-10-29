package db

import "context"

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User `json:"user"`
}

// CreateUserTx は口座から口座への送金を行う

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (
	CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)

		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)

	})

	return result, err
}
