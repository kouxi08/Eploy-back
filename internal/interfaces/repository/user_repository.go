package repository

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	FindUserByUID(ctx context.Context, tx *sql.Tx, uid string) (int, error)
	CreateUser(ctx context.Context, tx *sql.Tx, uid string) (int, error)
	FindOrCreateUserByUID(ctx context.Context, uid string) (int, error)
}
