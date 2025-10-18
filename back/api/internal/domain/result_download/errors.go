package result_download

import "errors"

var (
	// ErrInvalidResultID result ID is invalid
	ErrInvalidResultID = errors.New("結果IDが無効です")

	// ErrInvalidUserID user ID is invalid
	ErrInvalidUserID = errors.New("ユーザーIDが無効です")

	// ErrDownloadNotFound download record not found
	ErrDownloadNotFound = errors.New("ダウンロード履歴が見つかりません")

	// ErrDownloadAlreadyExists download record already exists
	ErrDownloadAlreadyExists = errors.New("既にダウンロード済みです")
)
