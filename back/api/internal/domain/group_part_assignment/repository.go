package group_part_assignment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// グループパーツ割り当てリポジトリのインターフェース
type Repository interface {
	// 新しいグループパーツ割り当てを作成
	Create(ctx context.Context, assignment *GroupPartAssignment) error

	// 割り当てIDでグループパーツ割り当てを検索
	FindByID(ctx context.Context, assignmentID uuid.UUID) (*GroupPartAssignment, error)

	// グループIDとコラージュ日で割り当てを取得
	FindByGroupAndDay(ctx context.Context, groupID string, collageDay time.Time) ([]*GroupPartAssignment, error)

	// ユーザーID、グループID、コラージュ日で割り当てを検索
	FindByUserGroupAndDay(ctx context.Context, userID uuid.UUID, groupID string, collageDay time.Time) (*GroupPartAssignment, error)

	// パーツIDで割り当てを検索
	FindByPartID(ctx context.Context, partID uuid.UUID) ([]*GroupPartAssignment, error)

	// グループパーツ割り当て情報を更新
	Update(ctx context.Context, assignment *GroupPartAssignment) error

	// グループパーツ割り当てを削除
	Delete(ctx context.Context, assignmentID uuid.UUID) error

	// グループIDとコラージュ日で全ての割り当てを削除
	DeleteByGroupAndDay(ctx context.Context, groupID string, collageDay time.Time) error

	// 全グループパーツ割り当てを取得
	List(ctx context.Context, limit, offset int) ([]*GroupPartAssignment, error)
}
