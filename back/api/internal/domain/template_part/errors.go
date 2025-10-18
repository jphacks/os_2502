package template_part

import "errors"

var (
	// ErrInvalidPartID パーツIDが無効
	ErrInvalidPartID = errors.New("パーツIDが無効です")

	// ErrInvalidTemplateID テンプレートIDが無効
	ErrInvalidTemplateID = errors.New("テンプレートIDが無効です")

	// ErrInvalidPartNumber パーツ番号が無効
	ErrInvalidPartNumber = errors.New("パーツ番号は1以上である必要があります")

	// ErrInvalidDimensions サイズが無効
	ErrInvalidDimensions = errors.New("幅と高さは0より大きい必要があります")

	// ErrTemplatePartNotFound テンプレートパーツが見つからない
	ErrTemplatePartNotFound = errors.New("テンプレートパーツが見つかりません")

	// ErrTemplatePartAlreadyExists テンプレートパーツが既に存在する
	ErrTemplatePartAlreadyExists = errors.New("テンプレートパーツは既に存在します")

	// ErrDuplicatePartNumber 同じテンプレート内でパーツ番号が重複
	ErrDuplicatePartNumber = errors.New("同じテンプレート内でパーツ番号が重複しています")
)
