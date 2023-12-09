package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nickyrolly/ws-chat-demo/internal/domain"
	"github.com/nickyrolly/ws-chat-demo/internal/repository/postgre"
	"github.com/nickyrolly/ws-chat-demo/pkg/auth"
)

func LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := postgre.FindUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	err = user.ComparePassword(password)
	if err != nil {
		return "", err
	}

	// Create a JWT
	token, err := auth.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func RegisterUser(ctx context.Context, username, password string) error {
	newUser := domain.NewUser(username, password)
	err := newUser.HashPassword()
	if err != nil {
		return err
	}

	existingUser, err := postgre.FindUserByUsername(ctx, username)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingUser.ID > 0 {
		return errors.New("user already exist")
	}

	err = postgre.SaveUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}
