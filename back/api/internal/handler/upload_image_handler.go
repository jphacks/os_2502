package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/upload_image"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type UploadImageHandler struct {
	useCase *usecase.UploadImageUseCase
}

func NewUploadImageHandler(useCase *usecase.UploadImageUseCase) *UploadImageHandler {
	return &UploadImageHandler{useCase: useCase}
}

type UploadImageRequest struct {
	FileURL    string `json:"file_url"`
	GroupID    string `json:"group_id"`
	CollageDay string `json:"collage_day"` // YYYY-MM-DD format
}

type UploadImageResponse struct {
	ImageID    string `json:"image_id"`
	FileURL    string `json:"file_url"`
	GroupID    string `json:"group_id"`
	UserID     string `json:"user_id"`
	CollageDay string `json:"collage_day"`
	CreatedAt  string `json:"created_at"`
}

func toUploadImageResponse(ui *upload_image.UploadImage) UploadImageResponse {
	return UploadImageResponse{
		ImageID:    ui.ImageID().String(),
		FileURL:    ui.FileURL(),
		GroupID:    ui.GroupID(),
		UserID:     ui.UserID().String(),
		CollageDay: ui.CollageDay().Format("2006-01-02"),
		CreatedAt:  ui.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *UploadImageHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
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

	var req UploadImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	collageDay, err := time.Parse("2006-01-02", req.CollageDay)
	if err != nil {
		respondError(w, http.StatusBadRequest, "コラージュ日が無効です（YYYY-MM-DD形式で指定してください）")
		return
	}

	image, err := h.useCase.UploadImage(r.Context(), req.FileURL, req.GroupID, userID, collageDay)
	if err != nil {
		switch err {
		case upload_image.ErrInvalidFileURL, upload_image.ErrInvalidGroupID, upload_image.ErrInvalidUserID:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "画像のアップロードに失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toUploadImageResponse(image))
}

func (h *UploadImageHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	idStr := r.URL.Path[len("/api/images/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	image, err := h.useCase.GetImage(r.Context(), id)
	if err != nil {
		if err == upload_image.ErrImageNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "画像の取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toUploadImageResponse(image))
}

func (h *UploadImageHandler) GetImagesByGroup(w http.ResponseWriter, r *http.Request) {
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

	images, err := h.useCase.GetImagesByGroup(r.Context(), groupID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "画像一覧の取得に失敗しました")
		return
	}

	responses := make([]UploadImageResponse, 0, len(images))
	for _, img := range images {
		responses = append(responses, toUploadImageResponse(img))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"images": responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

func (h *UploadImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

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

	idStr := r.URL.Path[len("/api/images/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "画像IDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効な画像IDです")
		return
	}

	if err := h.useCase.DeleteImage(r.Context(), id, userID); err != nil {
		switch err {
		case upload_image.ErrImageNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case upload_image.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "画像の削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
