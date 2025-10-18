package collage_result

import "errors"

var (
	// ErrInvalidTemplateID template ID is invalid
	ErrInvalidTemplateID = errors.New("テンプレートIDが無効です")

	// ErrInvalidGroupID group ID is invalid
	ErrInvalidGroupID = errors.New("グループIDが無効です")

	// ErrInvalidFileURL file URL is invalid
	ErrInvalidFileURL = errors.New("ファイルURLが無効です（1〜500文字で指定してください）")

	// ErrInvalidTargetUserNumber target user number is invalid
	ErrInvalidTargetUserNumber = errors.New("対象ユーザー数が無効です（1以上で指定してください）")

	// ErrResultNotFound result not found
	ErrResultNotFound = errors.New("コラージュ結果が見つかりません")

	// ErrResultAlreadyExists result already exists
	ErrResultAlreadyExists = errors.New("このコラージュ結果は既に存在します")
)
