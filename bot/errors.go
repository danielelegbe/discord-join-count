package bot

import (
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound = errors.New("User not found")
)

func HandleUserExists(err error) (bool, error) {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrUserNotFound
		}

		return false, err
	}
	return true, nil
}
