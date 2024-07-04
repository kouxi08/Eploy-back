package persistence

import (
	"context"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindUserByUID(ctx context.Context, tx *sql.Tx, uid string) (int, error) {
	var userID int
	query := `
		SELECT 
			id 
		FROM
			users 
		WHERE 
			external_uid = ? FOR UPDATE`
	err := tx.QueryRowContext(ctx, query, uid).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, tx *sql.Tx, uid string) (int, error) {
	query := `
		INSERT INTO 
			users 
				(external_uid)
		VALUES 
			(?)`
	result, err := tx.ExecContext(ctx, query, uid)
	if err != nil {
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(userID), nil
}

func (r *UserRepository) FindOrCreateUserByUID(ctx context.Context, uid string) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	userID, err := r.FindUserByUID(ctx, tx, uid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if userID == 0 {
		userID, err = r.CreateUser(ctx, tx, uid)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to create user: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userID, nil
}
