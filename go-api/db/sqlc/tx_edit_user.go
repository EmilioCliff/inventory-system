package db

import "context"

type EditUserParams struct {
	UserID      int64  `json:"user_id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Username    string `json:"username"`
}

type EditUserResult struct {
	UserEdited User `json:"useredited"`
}

func (store *Store) EditUserTx(ctx context.Context, arg EditUserParams) (EditUserResult, error) {
	var result EditUserResult

	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUserForUpdate(ctx, arg.UserID)
		if err != nil {
			return err
		}

		if arg.Role == "admin" {
			result.UserEdited, err = q.UpdateUserCredentials(ctx, UpdateUserCredentialsParams{
				Password:    user.Password,
				UserID:      user.UserID,
				Email:       arg.Email,
				Username:    arg.Username,
				Address:     arg.Address,
				PhoneNumber: arg.PhoneNumber,
			})
			if err != nil {
				return err
			}
		} else {
			result.UserEdited, err = q.UpdateUserCredentials(ctx, UpdateUserCredentialsParams{
				Password:    arg.Password,
				Email:       user.Email,
				UserID:      user.UserID,
				Username:    user.Username,
				Address:     user.Address,
				PhoneNumber: user.PhoneNumber,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}
