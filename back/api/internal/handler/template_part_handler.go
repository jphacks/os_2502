package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/template_part"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type TemplatePartHandler struct {
	useCase *usecase.TemplatePartUseCase
}

func NewTemplatePartHandler(useCase *usecase.TemplatePartUseCase) *TemplatePartHandler {
	return &TemplatePartHandler{useCase: useCase}
}

type CreateTemplatePartRequest struct {
	TemplateID  string  `json:"template_id"`
	PartNumber  int     `json:"part_number"`
	PositionX   int     `json:"position_x"`
	PositionY   int     `json:"position_y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	PartName    *string `json:"part_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateTemplatePartPositionRequest struct {
	PositionX int `json:"position_x"`
	PositionY int `json:"position_y"`
	Width     int `json:"width"`
	Height    int `json:"height"`
}

type UpdateTemplatePartNameRequest struct {
	PartName *string `json:"part_name"`
}

type UpdateTemplatePartDescriptionRequest struct {
	Description *string `json:"description"`
}

type TemplatePartResponse struct {
	PartID      string  `json:"part_id"`
	TemplateID  string  `json:"template_id"`
	PartNumber  int     `json:"part_number"`
	PartName    *string `json:"part_name,omitempty"`
	PositionX   int     `json:"position_x"`
	PositionY   int     `json:"position_y"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Description *string `json:"description,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// TemplatePartエンティティをTemplatePartResponseに変換
func toTemplatePartResponse(tp *template_part.TemplatePart) TemplatePartResponse {
	return TemplatePartResponse{
		PartID:      tp.PartID().String(),
		TemplateID:  tp.TemplateID().String(),
		PartNumber:  tp.PartNumber(),
		PartName:    tp.PartName(),
		PositionX:   tp.PositionX(),
		PositionY:   tp.PositionY(),
		Width:       tp.Width(),
		Height:      tp.Height(),
		Description: tp.Description(),
		CreatedAt:   tp.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   tp.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *TemplatePartHandler) CreateTemplatePart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	var req CreateTemplatePartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートIDです")
		return
	}

	tp, err := h.useCase.CreateTemplatePart(
		r.Context(),
		templateID,
		req.PartNumber,
		req.PositionX,
		req.PositionY,
		req.Width,
		req.Height,
		req.PartName,
		req.Description,
	)
	if err != nil {
		switch err {
		case template_part.ErrInvalidTemplateID, template_part.ErrInvalidPartNumber, template_part.ErrInvalidDimensions:
			respondError(w, http.StatusBadRequest, err.Error())
		case template_part.ErrDuplicatePartNumber, template_part.ErrTemplatePartAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "テンプレートパーツの作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toTemplatePartResponse(tp))
}

func (h *TemplatePartHandler) GetTemplatePart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/template-parts/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "テンプレートパーツIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートパーツIDです")
		return
	}

	tp, err := h.useCase.GetTemplatePartByID(r.Context(), id)
	if err != nil {
		if err == template_part.ErrTemplatePartNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "テンプレートパーツの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toTemplatePartResponse(tp))
}

func (h *TemplatePartHandler) UpdateTemplatePartPosition(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/template-parts/"):]
	if idStr == "" || len(idStr) < len("/position")+1 {
		respondError(w, http.StatusBadRequest, "テンプレートパーツIDが必要です")
		return
	}
	// "/position" を削除
	idStr = idStr[:len(idStr)-len("/position")]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートパーツIDです")
		return
	}

	var req UpdateTemplatePartPositionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	tp, err := h.useCase.UpdateTemplatePartPosition(r.Context(), id, req.PositionX, req.PositionY, req.Width, req.Height)
	if err != nil {
		switch err {
		case template_part.ErrInvalidDimensions:
			respondError(w, http.StatusBadRequest, err.Error())
		case template_part.ErrTemplatePartNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "テンプレートパーツの位置更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toTemplatePartResponse(tp))
}

func (h *TemplatePartHandler) UpdateTemplatePartName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/template-parts/"):]
	if idStr == "" || len(idStr) < len("/name")+1 {
		respondError(w, http.StatusBadRequest, "テンプレートパーツIDが必要です")
		return
	}
	// "/name" を削除
	idStr = idStr[:len(idStr)-len("/name")]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートパーツIDです")
		return
	}

	var req UpdateTemplatePartNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	tp, err := h.useCase.UpdateTemplatePartName(r.Context(), id, req.PartName)
	if err != nil {
		if err == template_part.ErrTemplatePartNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "テンプレートパーツ名の更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toTemplatePartResponse(tp))
}

func (h *TemplatePartHandler) UpdateTemplatePartDescription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/template-parts/"):]
	if idStr == "" || len(idStr) < len("/description")+1 {
		respondError(w, http.StatusBadRequest, "テンプレートパーツIDが必要です")
		return
	}
	// "/description" を削除
	idStr = idStr[:len(idStr)-len("/description")]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートパーツIDです")
		return
	}

	var req UpdateTemplatePartDescriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	tp, err := h.useCase.UpdateTemplatePartDescription(r.Context(), id, req.Description)
	if err != nil {
		if err == template_part.ErrTemplatePartNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "テンプレートパーツ説明の更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toTemplatePartResponse(tp))
}

func (h *TemplatePartHandler) DeleteTemplatePart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/template-parts/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "テンプレートパーツIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートパーツIDです")
		return
	}

	if err := h.useCase.DeleteTemplatePart(r.Context(), id); err != nil {
		if err == template_part.ErrTemplatePartNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "テンプレートパーツの削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TemplatePartHandler) ListTemplateParts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// クエリパラメータからlimitとoffsetを取得
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
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

	parts, err := h.useCase.ListTemplateParts(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレートパーツ一覧の取得に失敗しました")
		return
	}

	responses := make([]TemplatePartResponse, 0, len(parts))
	for _, tp := range parts {
		responses = append(responses, toTemplatePartResponse(tp))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"template_parts": responses,
		"limit":          limit,
		"offset":         offset,
		"count":          len(responses),
	})
}

func (h *TemplatePartHandler) GetTemplatePartsByTemplateID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	templateIDStr := r.URL.Query().Get("template_id")
	if templateIDStr == "" {
		respondError(w, http.StatusBadRequest, "テンプレートIDが必要です")
		return
	}

	templateID, err := uuid.Parse(templateIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートIDです")
		return
	}

	parts, err := h.useCase.GetTemplatePartsByTemplateID(r.Context(), templateID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "テンプレートパーツの取得に失敗しました")
		return
	}

	responses := make([]TemplatePartResponse, 0, len(parts))
	for _, tp := range parts {
		responses = append(responses, toTemplatePartResponse(tp))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"template_parts": responses,
		"count":          len(responses),
	})
}
