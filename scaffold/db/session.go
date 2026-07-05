package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateSession(ctx context.Context, userID int) (string, error) {
	token, err := randomHex(32)
	if err != nil {
		return "", err
	}
	if _, err := s.Pool.Exec(
		ctx,
		`INSERT INTO sessions(token, user_id, expires_at)
		 VALUES($1, $2, now() + interval '7 days')`,
		token,
		userID,
	); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Store) CurrentUser(ctx context.Context, token string) (User, bool, error) {
	var user User
	err := s.Pool.QueryRow(
		ctx,
		`SELECT users.id, users.email
		 FROM sessions JOIN users ON users.id = sessions.user_id
		 WHERE sessions.token = $1 AND sessions.expires_at > now()`,
		token,
	).Scan(&user.ID, &user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return user, true, nil
}

func (s *Store) DestroySession(ctx context.Context, token string) error {
	_, err := s.Pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}
