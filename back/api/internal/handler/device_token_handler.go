package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/device_token"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type DeviceTokenHandler struct {
	useCase *usecase.DeviceTokenUseCase
}

func NewDeviceTokenHandler(useCase *usecase.DeviceTokenUseCase) *DeviceTokenHandler {
	return &DeviceTokenHandler{useCase: useCase}
}

type RegisterDeviceTokenRequest struct {
	DeviceToken string  `json:"device_token"`
	DeviceType  string  `json:"device_type"`
	DeviceName  *string `json:"device_name,omitempty"`
}

type DeviceTokenResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	DeviceToken string  `json:"device_token"`
	DeviceType  string  `json:"device_type"`
	DeviceName  *string `json:"device_name,omitempty"`
	IsActive    bool    `json:"is_active"`
	LastUsedAt  *string `json:"last_used_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func toDeviceTokenResponse(dt *device_token.DeviceToken) DeviceTokenResponse {
	resp := DeviceTokenResponse{
		ID:          dt.ID().String(),
		UserID:      dt.UserID().String(),
		DeviceToken: dt.DeviceToken(),
		DeviceType:  string(dt.DeviceType()),
		DeviceName:  dt.DeviceName(),
		IsActive:    dt.IsActive(),
		CreatedAt:   dt.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   dt.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}

	if lastUsedAt := dt.LastUsedAt(); lastUsedAt != nil {
		formatted := lastUsedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.LastUsedAt = &formatted
	}

	return resp
}

func (h *DeviceTokenHandler) RegisterDeviceToken(w http.ResponseWriter, r *http.Request) {
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

	var req RegisterDeviceTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	deviceType := device_token.DeviceType(req.DeviceType)
	token, err := h.useCase.RegisterDeviceToken(r.Context(), userID, req.DeviceToken, deviceType, req.DeviceName)
	if err != nil {
		switch err {
		case device_token.ErrInvalidDeviceToken, device_token.ErrInvalidDeviceType:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "デバイストークンの登録に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toDeviceTokenResponse(token))
}

func (h *DeviceTokenHandler) GetDeviceToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/device-tokens/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "デバイストークンIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なデバイストークンIDです")
		return
	}

	token, err := h.useCase.GetDeviceToken(r.Context(), id)
	if err != nil {
		if err == device_token.ErrDeviceTokenNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "デバイストークンの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toDeviceTokenResponse(token))
}

func (h *DeviceTokenHandler) GetUserDeviceTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	tokens, err := h.useCase.GetUserDeviceTokens(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "デバイストークン一覧の取得に失敗しました")
		return
	}

	responses := make([]DeviceTokenResponse, 0, len(tokens))
	for _, t := range tokens {
		responses = append(responses, toDeviceTokenResponse(t))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"device_tokens": responses,
		"limit":         limit,
		"offset":        offset,
		"count":         len(responses),
	})
}

func (h *DeviceTokenHandler) DeactivateDeviceToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
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

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/device-tokens/"):]
	if idStr == "" || len(idStr) < len("/deactivate")+1 {
		respondError(w, http.StatusBadRequest, "デバイストークンIDが必要です")
		return
	}
	// "/deactivate" を削除
	idStr = idStr[:len(idStr)-len("/deactivate")]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なデバイストークンIDです")
		return
	}

	if err := h.useCase.DeactivateDeviceToken(r.Context(), id, userID); err != nil {
		if err == device_token.ErrDeviceTokenNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "デバイストークンの無効化に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DeviceTokenHandler) DeleteDeviceToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/device-tokens/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "デバイストークンIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なデバイストークンIDです")
		return
	}

	if err := h.useCase.DeleteDeviceToken(r.Context(), id, userID); err != nil {
		if err == device_token.ErrDeviceTokenNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "デバイストークンの削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
