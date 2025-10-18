package upload_images_collage_result

import (
	"context"

	"github.com/google/uuid"
)

// 画像とコラージュ結果の関連リポジトリのインターフェース
type Repository interface {
	// 新しい画像とコラージュ結果の関連を作成
	Create(ctx context.Context, relation *UploadImagesCollageResult) error

	// 画像IDとコラージュ結果IDで関連を検索
	FindByImageIDAndResultID(ctx context.Context, imageID, resultID uuid.UUID) (*UploadImagesCollageResult, error)

	// 画像IDで全ての関連を取得
	FindByImageID(ctx context.Context, imageID uuid.UUID) ([]*UploadImagesCollageResult, error)

	// コラージュ結果IDで全ての関連を取得
	FindByResultID(ctx context.Context, resultID uuid.UUID) ([]*UploadImagesCollageResult, error)

	// 画像とコラージュ結果の関連を更新
	Update(ctx context.Context, relation *UploadImagesCollageResult) error

	// 画像とコラージュ結果の関連を削除
	Delete(ctx context.Context, imageID, resultID uuid.UUID) error

	// コラージュ結果IDで全ての関連を削除
	DeleteByResultID(ctx context.Context, resultID uuid.UUID) error

	// 全ての関連を取得
	List(ctx context.Context, limit, offset int) ([]*UploadImagesCollageResult, error)
}
