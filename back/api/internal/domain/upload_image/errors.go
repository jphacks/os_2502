package upload_image

import "errors"

var (
	// ErrInvalidFileURL file URL is invalid
	ErrInvalidFileURL = errors.New("ファイルURLが無効です（1〜500文字で指定してください）")

	// ErrInvalidGroupID group ID is invalid
	ErrInvalidGroupID = errors.New("グループIDが無効です")

	// ErrInvalidUserID user ID is invalid
	ErrInvalidUserID = errors.New("ユーザーIDが無効です")

	// ErrImageNotFound image not found
	ErrImageNotFound = errors.New("画像が見つかりません")

	// ErrImageAlreadyExists image already exists
	ErrImageAlreadyExists = errors.New("この画像は既に存在します")

	// ErrNotAuthorized not authorized to access this image
	ErrNotAuthorized = errors.New("この画像にアクセスする権限がありません")
)
