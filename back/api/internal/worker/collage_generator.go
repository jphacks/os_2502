package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png" // PNGデコーダーを登録
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jphacks/os_2502/back/api/internal/domain/group"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_member"
)

// TemplateFrame テンプレートのフレーム情報
type TemplateFrame struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"width"`
	H    int    `json:"height"`
}

// TemplateData テンプレート情報
type TemplateData struct {
	Name       string          `json:"name"`
	PhotoCount int             `json:"photo_count"`
	ViewBox    string          `json:"viewBox"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
	Frames     []TemplateFrame `json:"frames"`
}

// CollageGenerator コラージュ生成ワーカー
type CollageGenerator struct {
	groupRepo       group.Repository
	groupMemberRepo group_member.Repository
	checkInterval   time.Duration
	templatesPath   string
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
		templatesPath:   "resources/templates.json",
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
	log.Printf("Generating collage for group %s from %s", groupID, uploadDir)

	// グループ情報を取得してテンプレートIDを確認
	g, err := w.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	templateID := g.TemplateID()
	if templateID == nil || *templateID == "" {
		return fmt.Errorf("template ID not found in group")
	}

	log.Printf("Using template ID: %s", *templateID)

	// テンプレート情報を読み込み
	template, err := w.loadTemplate(*templateID)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	log.Printf("Loaded template: %s (%dx%d)", template.Name, template.Width, template.Height)

	// アップロードされた画像ファイルを取得
	files, err := os.ReadDir(uploadDir)
	if err != nil {
		return fmt.Errorf("failed to read upload directory: %w", err)
	}

	// 画像ファイルのパスリストを作成
	var imagePaths []string
	for _, file := range files {
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				imagePaths = append(imagePaths, filepath.Join(uploadDir, file.Name()))
			}
		}
	}

	if len(imagePaths) != template.PhotoCount {
		return fmt.Errorf("image count mismatch: expected %d, got %d", template.PhotoCount, len(imagePaths))
	}

	// コラージュ画像を生成
	resultImage, err := w.createCollageImage(template, imagePaths)
	if err != nil {
		return fmt.Errorf("failed to create collage image: %w", err)
	}

	// 結果ディレクトリを作成
	resultDir := "/uploads/collages"
	if err := os.MkdirAll(resultDir, 0755); err != nil {
		return fmt.Errorf("failed to create result directory: %w", err)
	}

	// コラージュ画像を保存
	resultPath := filepath.Join(resultDir, groupID+"_collage.jpg")
	outFile, err := os.Create(resultPath)
	if err != nil {
		return fmt.Errorf("failed to create result file: %w", err)
	}
	defer outFile.Close()

	if err := jpeg.Encode(outFile, resultImage, &jpeg.Options{Quality: 90}); err != nil {
		return fmt.Errorf("failed to encode collage image: %w", err)
	}

	log.Printf("Collage saved to %s", resultPath)

	return nil
}

// loadTemplate テンプレート情報を読み込み
func (w *CollageGenerator) loadTemplate(templateID string) (*TemplateData, error) {
	// templates.jsonファイルを読み込む
	filePath := w.templatesPath
	if !filepath.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
		filePath = filepath.Join(wd, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates file: %w", err)
	}

	var templates []TemplateData
	if err := json.Unmarshal(data, &templates); err != nil {
		return nil, fmt.Errorf("failed to parse templates data: %w", err)
	}

	// テンプレートIDで検索（現在はnameで検索）
	for _, t := range templates {
		if t.Name == templateID {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("template not found: %s", templateID)
}

// createCollageImage コラージュ画像を作成
func (w *CollageGenerator) createCollageImage(template *TemplateData, imagePaths []string) (image.Image, error) {
	// キャンバスを作成（デフォルトサイズ: 1000x1000）
	width := template.Width
	height := template.Height
	if width == 0 {
		width = 1000
	}
	if height == 0 {
		height = 1000
	}

	canvas := image.NewRGBA(image.Rect(0, 0, width, height))

	// 各フレームに画像を配置
	for i, frame := range template.Frames {
		if i >= len(imagePaths) {
			break
		}

		// 画像を読み込み
		imgFile, err := os.Open(imagePaths[i])
		if err != nil {
			log.Printf("Warning: failed to open image %s: %v", imagePaths[i], err)
			continue
		}

		img, _, err := image.Decode(imgFile)
		imgFile.Close()
		if err != nil {
			log.Printf("Warning: failed to decode image %s: %v", imagePaths[i], err)
			continue
		}

		// フレームの位置とサイズ
		x, y, fw, fh := frame.X, frame.Y, frame.W, frame.H
		if fw == 0 || fh == 0 {
			// サイズが指定されていない場合は画像のサイズを使用
			fw = img.Bounds().Dx()
			fh = img.Bounds().Dy()
		}

		// 画像をリサイズして配置
		resized := w.resizeImage(img, fw, fh)
		dp := image.Point{X: x, Y: y}
		dr := image.Rectangle{Min: dp, Max: dp.Add(resized.Bounds().Size())}
		draw.Draw(canvas, dr, resized, image.Point{}, draw.Over)

		log.Printf("Placed image %d at (%d,%d) size (%dx%d)", i, x, y, fw, fh)
	}

	return canvas, nil
}

// resizeImage 画像をリサイズ（簡易実装）
func (w *CollageGenerator) resizeImage(img image.Image, width, height int) image.Image {
	// 簡易的なニアレストネイバー法でリサイズ
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := x * srcW / width
			srcY := y * srcH / height
			dst.Set(x, y, img.At(bounds.Min.X+srcX, bounds.Min.Y+srcY))
		}
	}

	return dst
}
