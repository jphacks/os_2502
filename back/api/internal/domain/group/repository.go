package group

import "context"

// Repository はグループのリポジトリインターフェース
type Repository interface {
	// Create は新しいグループを作成
	Create(ctx context.Context, group *Group) error

	// FindByID はIDでグループを検索
	FindByID(ctx context.Context, id string) (*Group, error)

	// FindByInvitationToken は招待トークンでグループを検索
	FindByInvitationToken(ctx context.Context, token string) (*Group, error)

	// FindByOwnerUserID はオーナーユーザーIDでグループを検索
	FindByOwnerUserID(ctx context.Context, ownerUserID string, limit, offset int) ([]*Group, error)

	// List は全グループを取得
	List(ctx context.Context, limit, offset int) ([]*Group, error)

	// Update はグループ情報を更新
	Update(ctx context.Context, group *Group) error

	// Delete はグループを削除
	Delete(ctx context.Context, id string) error

	// Count は全グループ数を取得
	Count(ctx context.Context) (int, error)

	// CountByOwnerUserID は特定オーナーのグループ数を取得
	CountByOwnerUserID(ctx context.Context, ownerUserID string) (int, error)

	// FindByStatus はステータスでグループを検索
	FindByStatus(ctx context.Context, status string, limit, offset int) ([]*Group, error)

	// UpdateStatus はグループのステータスを更新
	UpdateStatus(ctx context.Context, id string, status string) error
}
