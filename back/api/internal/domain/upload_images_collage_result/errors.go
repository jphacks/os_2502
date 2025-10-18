package upload_images_collage_result

import "errors"

var (
	// ErrInvalidImageID 画像IDが無効
	ErrInvalidImageID = errors.New("画像IDが無効です")

	// ErrInvalidResultID 結果IDが無効
	ErrInvalidResultID = errors.New("結果IDが無効です")

	// ErrInvalidDimensions サイズが無効
	ErrInvalidDimensions = errors.New("幅と高さは0より大きい必要があります")

	// ErrUploadImagesCollageResultNotFound 画像とコラージュ結果の関連が見つからない
	ErrUploadImagesCollageResultNotFound = errors.New("画像とコラージュ結果の関連が見つかりません")

	// ErrUploadImagesCollageResultAlreadyExists 画像とコラージュ結果の関連が既に存在する
	ErrUploadImagesCollageResultAlreadyExists = errors.New("画像とコラージュ結果の関連は既に存在します")
)
