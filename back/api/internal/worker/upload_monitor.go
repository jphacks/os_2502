package worker

import (
	"context"

	"github.com/jphacks/os_2502/back/api/internal/domain/upload_image"
)

// UploadStatus アップロード状況
type UploadStatus struct {
	GroupID  string `json:"group_id"`
	Total    int    `json:"total"`
	Uploaded int    `json:"uploaded"`
	Status   string `json:"status"`
}

// UploadMonitor アップロード監視
type UploadMonitor struct {
	uploadImageRepo upload_image.Repository
}

// NewUploadMonitor UploadMonitorを作成
func NewUploadMonitor(uploadImageRepo upload_image.Repository) *UploadMonitor {
	return &UploadMonitor{
		uploadImageRepo: uploadImageRepo,
	}
}

// CheckUploadStatus アップロード状況を確認
func (m *UploadMonitor) CheckUploadStatus(ctx context.Context, groupID string) (*UploadStatus, error) {
	// グループのアップロード画像を取得
	images, err := m.uploadImageRepo.FindByGroupID(ctx, groupID, 1000, 0)
	if err != nil {
		return nil, err
	}

	status := &UploadStatus{
		GroupID:  groupID,
		Total:    0,
		Uploaded: len(images),
		Status:   "in_progress",
	}

	// TODO: グループメンバー数を取得して Total に設定
	// 現在は仮の値を設定
	status.Total = 4

	// 全員アップロード完了の場合
	if status.Uploaded >= status.Total {
		status.Status = "completed"
	}

	return status, nil
}
