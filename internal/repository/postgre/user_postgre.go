package postgre

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/nickyrolly/ws-chat-demo/internal/domain"
)

var (
	// SELECT QUERY
	QuerySelectUserByUsername string = `SELECT id, username, password FROM public.user WHERE username = $1`
	QuerySelectUserByID       string = `SELECT id, username, password FROM public.user WHERE id = $1`

	// INSERT QUERY
	QueryInsertUser string = `INSERT INTO public.user (username, password) VALUES ($1, $2)`
)

func FindUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User

	row := DBChat.QueryRowContext(ctx, QuerySelectUserByUsername, username)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func FindUserByID(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User

	row := DBChat.QueryRowContext(ctx, QuerySelectUserByID, id)
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Println("Error: Can't find user")
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}

func SaveUser(ctx context.Context, user domain.User) error {
	_, err := DBChat.ExecContext(ctx, QueryInsertUser, user.Username, user.Password)

	if err != nil {
		return err
	}

	return nil
}
