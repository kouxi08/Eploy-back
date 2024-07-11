package persistence

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestFindUserByUID(t *testing.T) {
	repo := NewUserRepository(db)
	ctx := context.Background()

	// テスト用のトランザクションを開始
	tx, err := db.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	// テスト用のユーザーを作成
	uid := "test_uid_1"
	sqlResult, err := tx.Exec("INSERT INTO users (external_uid) VALUES (?)", uid)
	assert.NoError(t, err)
	wantId, err := sqlResult.LastInsertId()
	assert.NoError(t, err)

	// ユーザーの検索
	userID, err := repo.FindUserByUID(ctx, tx, uid)

	// 検索結果の検証
	assert.Equal(t, int(wantId), userID, "UserID should match the created user")

	assert.NotEqual(t, 0, userID, "UserID should not be 0")
}

func TestCreateUser(t *testing.T) {
	repo := NewUserRepository(db)
	ctx := context.Background()

	// テスト用のトランザクションを開始
	tx, err := db.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	// ユーザーの作成
	uid := "test_uid_2"
	userID, err := repo.CreateUser(ctx, tx, uid)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userID, "UserID should not be 0")

	// 作成したユーザーの検証
	var storedUID string
	err = tx.QueryRow("SELECT external_uid FROM users WHERE id = ?", userID).Scan(&storedUID)
	assert.NoError(t, err)
	assert.Equal(t, uid, storedUID, "Stored UID should match the created UID")
}

func TestFindOrCreateUserByUID(t *testing.T) {
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 既存のユーザーを検索
	uidExisting := "test_uid_existing"
	tx, err := db.BeginTx(ctx, nil)
	assert.NoError(t, err)
	_, err = tx.Exec("INSERT INTO users (external_uid) VALUES (?)", uidExisting)
	assert.NoError(t, err)
	tx.Commit()

	userIDExisting, err := repo.FindOrCreateUserByUID(ctx, uidExisting)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userIDExisting, "UserID should not be 0")

	// 新規のユーザーを作成
	uidNew := "test_uid_new"
	userIDNew, err := repo.FindOrCreateUserByUID(ctx, uidNew)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, userIDNew, "UserID should not be 0")
	assert.NotEqual(t, userIDExisting, userIDNew, "UserID for new user should be different from existing user")
}
