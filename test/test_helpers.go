package test

import (
	"context"
	"database/sql"
)

func CreateUser(ctx context.Context, db *sql.DB) (int, error) {
	uid := "test_user"
	id, err := db.ExecContext(ctx, "INSERT INTO users (external_uid) VALUES (?)", uid)
	if err != nil {
		panic(err)
	}
	userId, err := id.LastInsertId()
	if err != nil {
		panic(err)
	}
	return int(userId), nil
}
