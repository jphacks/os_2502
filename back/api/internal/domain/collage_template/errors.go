package collage_template

import "errors"

var (
	// ErrInvalidName template name is invalid
	ErrInvalidName = errors.New("テンプレート名が無効です（1〜100文字で指定してください）")

	// ErrInvalidFilePath template file path is invalid
	ErrInvalidFilePath = errors.New("ファイルパスが無効です（1〜255文字で指定してください）")

	// ErrTemplateNotFound template not found
	ErrTemplateNotFound = errors.New("テンプレートが見つかりません")

	// ErrTemplateAlreadyExists template already exists
	ErrTemplateAlreadyExists = errors.New("このテンプレートは既に存在します")
)
