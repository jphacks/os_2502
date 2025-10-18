package template_part

import (
	"context"

	"github.com/google/uuid"
)

// テンプレートパーツリポジトリのインターフェース
type Repository interface {
	// 新しいテンプレートパーツを作成
	Create(ctx context.Context, part *TemplatePart) error

	// パーツIDでテンプレートパーツを検索
	FindByID(ctx context.Context, partID uuid.UUID) (*TemplatePart, error)

	// テンプレートIDで全てのパーツを取得
	FindByTemplateID(ctx context.Context, templateID uuid.UUID) ([]*TemplatePart, error)

	// テンプレートIDとパーツ番号でテンプレートパーツを検索
	FindByTemplateIDAndPartNumber(ctx context.Context, templateID uuid.UUID, partNumber int) (*TemplatePart, error)

	// テンプレートパーツ情報を更新
	Update(ctx context.Context, part *TemplatePart) error

	// テンプレートパーツを削除
	Delete(ctx context.Context, partID uuid.UUID) error

	// テンプレートIDに紐づく全てのパーツを削除
	DeleteByTemplateID(ctx context.Context, templateID uuid.UUID) error

	// 全テンプレートパーツを取得
	List(ctx context.Context, limit, offset int) ([]*TemplatePart, error)
}
