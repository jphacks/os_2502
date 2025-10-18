package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/collage_template"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type CollageTemplateHandler struct {
	useCase *usecase.CollageTemplateUseCase
}

func NewCollageTemplateHandler(useCase *usecase.CollageTemplateUseCase) *CollageTemplateHandler {
	return &CollageTemplateHandler{useCase: useCase}
}

type CreateTemplateRequest struct {
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

type UpdateTemplateRequest struct {
	Name     string `json:"name,omitempty"`
	FilePath string `json:"file_path,omitempty"`
}

type TemplateResponse struct {
	TemplateID string `json:"template_id"`
	Name       string `json:"name"`
	FilePath   string `json:"file_path"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

func toTemplateResponse(t *collage_template.CollageTemplate) TemplateResponse {
	return TemplateResponse{
		TemplateID: t.TemplateID().String(),
		Name:       t.Name(),
		FilePath:   t.FilePath(),
		CreatedAt:  t.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  t.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *CollageTemplateHandler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	template, err := h.useCase.CreateTemplate(r.Context(), req.Name, req.FilePath)
	if err != nil {
		switch err {
		case collage_template.ErrInvalidName, collage_template.ErrInvalidFilePath:
			respondError(w, http.StatusBadRequest, err.Error())
		case collage_template.ErrTemplateAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "テンプレートの作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toTemplateResponse(template))
}

func (h *CollageTemplateHandler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/templates/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "テンプレートIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートIDです")
		return
	}

	template, err := h.useCase.GetTemplate(r.Context(), id)
	if err != nil {
		if err == collage_template.ErrTemplateNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "テンプレートの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toTemplateResponse(template))
}

func (h *CollageTemplateHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	templates, err := h.useCase.ListTemplates(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレート一覧の取得に失敗しました")
		return
	}

	responses := make([]TemplateResponse, 0, len(templates))
	for _, t := range templates {
		responses = append(responses, toTemplateResponse(t))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"templates": responses,
		"limit":     limit,
		"offset":    offset,
		"count":     len(responses),
	})
}

func (h *CollageTemplateHandler) UpdateTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/templates/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "テンプレートIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートIDです")
		return
	}

	var req UpdateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	template, err := h.useCase.UpdateTemplate(r.Context(), id, req.Name, req.FilePath)
	if err != nil {
		switch err {
		case collage_template.ErrInvalidName, collage_template.ErrInvalidFilePath:
			respondError(w, http.StatusBadRequest, err.Error())
		case collage_template.ErrTemplateNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "テンプレートの更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toTemplateResponse(template))
}

func (h *CollageTemplateHandler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/templates/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "テンプレートIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートIDです")
		return
	}

	if err := h.useCase.DeleteTemplate(r.Context(), id); err != nil {
		if err == collage_template.ErrTemplateNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "テンプレートの削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
