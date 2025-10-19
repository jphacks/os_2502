package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// TemplateFrame represents a frame in a collage template
type TemplateFrame struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
}

// TemplateData represents a collage template with frames
type TemplateData struct {
	Name       string          `json:"name"`
	PhotoCount int             `json:"photo_count"`
	ViewBox    string          `json:"viewBox"`
	Frames     []TemplateFrame `json:"frames"`
}

type TemplateDataHandler struct {
	templatesPath string
}

func NewTemplateDataHandler() *TemplateDataHandler {
	return &TemplateDataHandler{
		templatesPath: "resources/templates.json",
	}
}

// GetTemplates returns all available collage templates
func (h *TemplateDataHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// templates.jsonファイルを読み込む
	filePath := h.templatesPath
	if !filepath.IsAbs(filePath) {
		// 相対パスの場合、プロジェクトルートからの相対パスとして解決
		wd, err := os.Getwd()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "ワーキングディレクトリの取得に失敗しました")
			return
		}
		filePath = filepath.Join(wd, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレートファイルの読み込みに失敗しました")
		return
	}

	// JSONをパース
	var templates []TemplateData
	if err := json.Unmarshal(data, &templates); err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレートデータの解析に失敗しました")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
	})
}

// GetTemplateByPhotoCount returns templates filtered by photo count
func (h *TemplateDataHandler) GetTemplateByPhotoCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// クエリパラメータからphoto_countを取得
	photoCountStr := r.URL.Query().Get("photo_count")
	if photoCountStr == "" {
		respondError(w, http.StatusBadRequest, "photo_countパラメータが必要です")
		return
	}

	var photoCount int
	if _, err := fmt.Sscanf(photoCountStr, "%d", &photoCount); err != nil {
		respondError(w, http.StatusBadRequest, "photo_countは数値である必要があります")
		return
	}

	// templates.jsonファイルを読み込む
	filePath := h.templatesPath
	if !filepath.IsAbs(filePath) {
		wd, err := os.Getwd()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "ワーキングディレクトリの取得に失敗しました")
			return
		}
		filePath = filepath.Join(wd, filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレートファイルの読み込みに失敗しました")
		return
	}

	// JSONをパース
	var allTemplates []TemplateData
	if err := json.Unmarshal(data, &allTemplates); err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレートデータの解析に失敗しました")
		return
	}

	// photo_countでフィルタリング
	var filteredTemplates []TemplateData
	for _, t := range allTemplates {
		if t.PhotoCount == photoCount {
			filteredTemplates = append(filteredTemplates, t)
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"templates":   filteredTemplates,
		"count":       len(filteredTemplates),
		"photo_count": photoCount,
	})
}
