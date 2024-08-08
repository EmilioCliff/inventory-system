package db

import "context"

func (store *Store) DeleteUserTx(ctx context.Context, user_id int32) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := store.DeleteUserStockValue(ctx, user_id)
		if err != nil {
			return err
		}

		err = store.DeleteUser(ctx, int64(user_id))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
