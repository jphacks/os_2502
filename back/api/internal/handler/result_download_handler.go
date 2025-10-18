package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/result_download"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type ResultDownloadHandler struct {
	useCase *usecase.ResultDownloadUseCase
}

func NewResultDownloadHandler(useCase *usecase.ResultDownloadUseCase) *ResultDownloadHandler {
	return &ResultDownloadHandler{useCase: useCase}
}

type RecordDownloadRequest struct {
	ResultID string `json:"result_id"`
}

type DownloadResponse struct {
	ResultID     string `json:"result_id"`
	UserID       string `json:"user_id"`
	DownloadedAt string `json:"downloaded_at"`
}

func toDownloadResponse(rd *result_download.ResultDownload) DownloadResponse {
	return DownloadResponse{
		ResultID:     rd.ResultID().String(),
		UserID:       rd.UserID().String(),
		DownloadedAt: rd.DownloadedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *ResultDownloadHandler) RecordDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	var req RecordDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	resultID, err := uuid.Parse(req.ResultID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	download, err := h.useCase.RecordDownload(r.Context(), resultID, userID)
	if err != nil {
		switch err {
		case result_download.ErrInvalidResultID, result_download.ErrInvalidUserID:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "ダウンロード記録の保存に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toDownloadResponse(download))
}

func (h *ResultDownloadHandler) GetDownloadsByResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
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

	downloads, err := h.useCase.GetDownloadsByResult(r.Context(), resultID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ダウンロード履歴の取得に失敗しました")
		return
	}

	responses := make([]DownloadResponse, 0, len(downloads))
	for _, dl := range downloads {
		responses = append(responses, toDownloadResponse(dl))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"downloads": responses,
		"limit":     limit,
		"offset":    offset,
		"count":     len(responses),
	})
}

func (h *ResultDownloadHandler) GetDownloadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	resultIDStr := r.URL.Query().Get("result_id")
	if resultIDStr == "" {
		respondError(w, http.StatusBadRequest, "結果IDが必要です")
		return
	}

	resultID, err := uuid.Parse(resultIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な結果IDです")
		return
	}

	count, err := h.useCase.GetDownloadCount(r.Context(), resultID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ダウンロード数の取得に失敗しました")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"result_id": resultIDStr,
		"count":     count,
	})
}
