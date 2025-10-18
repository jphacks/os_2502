package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/collage_result"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type CollageResultHandler struct {
	useCase *usecase.CollageResultUseCase
}

func NewCollageResultHandler(useCase *usecase.CollageResultUseCase) *CollageResultHandler {
	return &CollageResultHandler{useCase: useCase}
}

type CreateResultRequest struct {
	TemplateID       string `json:"template_id"`
	GroupID          string `json:"group_id"`
	FileURL          string `json:"file_url"`
	TargetUserNumber int    `json:"target_user_number"`
}

type CollageResultResponse struct {
	ResultID         string `json:"result_id"`
	TemplateID       string `json:"template_id"`
	GroupID          string `json:"group_id"`
	FileURL          string `json:"file_url"`
	TargetUserNumber int    `json:"target_user_number"`
	IsNotification   bool   `json:"is_notification"`
	CreatedAt        string `json:"created_at"`
}

func toCollageResultResponse(cr *collage_result.CollageResult) CollageResultResponse {
	return CollageResultResponse{
		ResultID:         cr.ResultID().String(),
		TemplateID:       cr.TemplateID().String(),
		GroupID:          cr.GroupID(),
		FileURL:          cr.FileURL(),
		TargetUserNumber: cr.TargetUserNumber(),
		IsNotification:   cr.IsNotification(),
		CreatedAt:        cr.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *CollageResultHandler) CreateResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	var req CreateResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なテンプレートIDです")
		return
	}

	result, err := h.useCase.CreateResult(r.Context(), templateID, req.GroupID, req.FileURL, req.TargetUserNumber)
	if err != nil {
		switch err {
		case collage_result.ErrInvalidTemplateID, collage_result.ErrInvalidGroupID, collage_result.ErrInvalidFileURL, collage_result.ErrInvalidTargetUserNumber:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "コラージュ結果の作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toCollageResultResponse(result))
}

func (h *CollageResultHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	idStr := r.URL.Path[len("/api/results/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	result, err := h.useCase.GetResult(r.Context(), id)
	if err != nil {
		if err == collage_result.ErrResultNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "コラージュ結果の取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toCollageResultResponse(result))
}

func (h *CollageResultHandler) GetResultsByGroup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
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

	results, err := h.useCase.GetResultsByGroup(r.Context(), groupID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "コラージュ結果一覧の取得に失敗しました")
		return
	}

	responses := make([]CollageResultResponse, 0, len(results))
	for _, res := range results {
		responses = append(responses, toCollageResultResponse(res))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"results": responses,
		"limit":   limit,
		"offset":  offset,
		"count":   len(responses),
	})
}

func (h *CollageResultHandler) MarkAsNotified(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	idStr := r.URL.Path[len("/api/results/"):]
	if idStr == "" || len(idStr) < len("/notify")+1 {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}
	// "/notify" を削除
	idStr = idStr[:len(idStr)-len("/notify")]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	if err := h.useCase.MarkAsNotified(r.Context(), id); err != nil {
		if err == collage_result.ErrResultNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "通知ステータスの更新に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CollageResultHandler) DeleteResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	idStr := r.URL.Path[len("/api/results/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	if err := h.useCase.DeleteResult(r.Context(), id); err != nil {
		if err == collage_result.ErrResultNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "コラージュ結果の削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
