package db

import "strings"

// MySQLの重複エラー（1062: Duplicate entry）
func IsDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// MySQL error code 1062: Duplicate entry
	return strings.Contains(errStr, "Error 1062") || strings.Contains(errStr, "Duplicate entry")
}

// MySQLの外部キー制約エラー（1451, 1452）
// エラーコード:
// - 1451: Cannot delete or update a parent row (参照されているレコードの削除)
// - 1452: Cannot add or update a child row (存在しない親レコードへの参照)
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Error 1451") ||
		strings.Contains(errStr, "Error 1452") ||
		strings.Contains(errStr, "foreign key constraint")
}
