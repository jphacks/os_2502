package user

import (
	"context"

	"github.com/google/uuid"
)

// ユーザーリポジトリのインターフェース
type Repository interface {
	// 新しいユーザーを作成
	Create(ctx context.Context, user *User) error

	// IDでユーザーを検索
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)

	// Firebase UIDでユーザーを検索
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*User, error)

	// usernameでユーザーを検索
	FindByUsername(ctx context.Context, username string) (*User, error)

	// usernameで部分一致検索
	SearchByUsername(ctx context.Context, query string, limit, offset int) ([]*User, error)

	// ユーザー情報を更新
	Update(ctx context.Context, user *User) error

	// ユーザーを削除
	Delete(ctx context.Context, id uuid.UUID) error

	// 全ユーザーを取得
	List(ctx context.Context, limit, offset int) ([]*User, error)
}
