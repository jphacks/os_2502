package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jphacks/os_2502/back/api/internal/domain/group"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_member"
)

// CollageGenerator コラージュ生成ワーカー
type CollageGenerator struct {
	groupRepo       group.Repository
	groupMemberRepo group_member.Repository
	checkInterval   time.Duration
}

// NewCollageGenerator コラージュ生成ワーカーを作成
func NewCollageGenerator(
	groupRepo group.Repository,
	groupMemberRepo group_member.Repository,
	checkInterval time.Duration,
) *CollageGenerator {
	if checkInterval == 0 {
		checkInterval = 10 * time.Second // デフォルト10秒
	}

	return &CollageGenerator{
		groupRepo:       groupRepo,
		groupMemberRepo: groupMemberRepo,
		checkInterval:   checkInterval,
	}
}

// Start ワーカーを開始（バックグラウンドで実行）
func (w *CollageGenerator) Start(ctx context.Context) {
	log.Println("🎨 Collage generator worker started")

	ticker := time.NewTicker(w.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🎨 Collage generator worker stopped")
			return
		case <-ticker.C:
			w.checkAndGenerateCollages(ctx)
		}
	}
}

// checkAndGenerateCollages カウントダウン状態のグループをチェックしてコラージュを生成
func (w *CollageGenerator) checkAndGenerateCollages(ctx context.Context) {
	// カウントダウン状態のグループを取得
	groups, err := w.groupRepo.FindByStatus(ctx, "countdown", 100, 0)
	if err != nil {
		log.Printf("❌ Failed to fetch countdown groups: %v", err)
		return
	}

	if len(groups) == 0 {
		return
	}

	log.Printf("🔍 Found %d groups in countdown status", len(groups))

	for _, g := range groups {
		if err := w.processGroup(ctx, g); err != nil {
			log.Printf("❌ Failed to process group %s: %v", g.ID(), err)
		}
	}
}

// processGroup グループの写真が全て揃っているかチェックし、コラージュを生成
func (w *CollageGenerator) processGroup(ctx context.Context, g *group.Group) error {
	groupID := g.ID()
	log.Printf("🔍 Checking group %s", groupID)

	// グループメンバーを取得
	members, err := w.groupMemberRepo.FindByGroupID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	memberCount := len(members)
	if memberCount == 0 {
		return fmt.Errorf("no members in group")
	}

	// アップロードされた写真をチェック
	uploadDir := "/uploads/groups/" + groupID
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		log.Printf("⏳ Upload directory not found for group %s", groupID)
		return nil
	}

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		return fmt.Errorf("failed to read upload directory: %w", err)
	}

	uploadedCount := 0
	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".jpg" || filepath.Ext(file.Name()) == ".jpeg" || filepath.Ext(file.Name()) == ".png") {
			uploadedCount++
		}
	}

	log.Printf("📊 Group %s: %d/%d photos uploaded", groupID, uploadedCount, memberCount)

	// 全員の写真が揃っていない場合はスキップ
	if uploadedCount < memberCount {
		return nil
	}

	log.Printf("✅ All photos uploaded for group %s, generating collage...", groupID)

	// コラージュを生成
	if err := w.generateCollage(ctx, groupID, uploadDir); err != nil {
		return fmt.Errorf("failed to generate collage: %w", err)
	}

	// グループステータスを完了に更新
	if err := w.groupRepo.UpdateStatus(ctx, groupID, "completed"); err != nil {
		log.Printf("⚠️ Failed to update group status: %v", err)
	}

	log.Printf("🎉 Collage generated successfully for group %s", groupID)

	// TODO: プッシュ通知を送信

	return nil
}

// generateCollage コラージュ画像を生成
func (w *CollageGenerator) generateCollage(ctx context.Context, groupID, uploadDir string) error {
	// TODO: 実際のコラージュ生成ロジックを実装
	// - テンプレートを取得
	// - 各メンバーの写真を読み込み
	// - テンプレートに従って配置
	// - 合成画像を保存

	log.Printf("🎨 Generating collage for group %s from %s", groupID, uploadDir)

	// 仮実装: コラージュ結果ディレクトリを作成
	resultDir := "/uploads/collages"
	if err := os.MkdirAll(resultDir, 0755); err != nil {
		return fmt.Errorf("failed to create result directory: %w", err)
	}

	// 仮実装: 空のコラージュファイルを作成（実際には画像合成処理）
	resultPath := filepath.Join(resultDir, groupID+"_collage.jpg")
	if err := os.WriteFile(resultPath, []byte("TODO: Generate actual collage"), 0644); err != nil {
		return fmt.Errorf("failed to write collage file: %w", err)
	}

	log.Printf("💾 Collage saved to %s", resultPath)

	return nil
}
