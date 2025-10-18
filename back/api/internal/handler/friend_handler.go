package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jphacks/os_2502/back/api/internal/domain/friend"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type FriendHandler struct {
	useCase *usecase.FriendUseCase
}

func NewFriendHandler(useCase *usecase.FriendUseCase) *FriendHandler {
	return &FriendHandler{useCase: useCase}
}

type SendFriendRequestRequest struct {
	AddresseeID string `json:"addressee_id"`
}

type FriendResponse struct {
	ID          string `json:"id"`
	RequesterID string `json:"requester_id"`
	AddresseeID string `json:"addressee_id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toFriendResponse(f *friend.Friend) FriendResponse {
	return FriendResponse{
		ID:          f.ID(),
		RequesterID: f.RequesterID(),
		AddresseeID: f.AddresseeID(),
		Status:      string(f.Status()),
		CreatedAt:   f.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   f.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *FriendHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	requesterID := r.Header.Get("X-User-ID")
	if requesterID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
		return
	}

	var req SendFriendRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	friendRequest, err := h.useCase.SendFriendRequest(r.Context(), requesterID, req.AddresseeID)
	if err != nil {
		switch err {
		case friend.ErrCannotFriendSelf, friend.ErrInvalidUserID:
			respondError(w, http.StatusBadRequest, err.Error())
		case friend.ErrAlreadyFriends, friend.ErrFriendRequestAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "フレンドリクエストの送信に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toFriendResponse(friendRequest))
}

func (h *FriendHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
		return
	}

	// URLパスからリクエストIDを取得
	requestID := r.URL.Path[len("/api/friends/"):]
	if requestID == "" || len(requestID) < len("/accept")+1 {
		respondError(w, http.StatusBadRequest, "リクエストIDが必要です")
		return
	}
	// "/accept" を削除
	requestID = requestID[:len(requestID)-len("/accept")]

	friendRequest, err := h.useCase.AcceptFriendRequest(r.Context(), requestID, userID)
	if err != nil {
		switch err {
		case friend.ErrFriendRequestNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case friend.ErrCannotAcceptNonPending:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "フレンドリクエストの承認に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toFriendResponse(friendRequest))
}

func (h *FriendHandler) RejectFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
		return
	}

	// URLパスからリクエストIDを取得
	requestID := r.URL.Path[len("/api/friends/"):]
	if requestID == "" || len(requestID) < len("/reject")+1 {
		respondError(w, http.StatusBadRequest, "リクエストIDが必要です")
		return
	}
	// "/reject" を削除
	requestID = requestID[:len(requestID)-len("/reject")]

	err := h.useCase.RejectFriendRequest(r.Context(), requestID, userID)
	if err != nil {
		switch err {
		case friend.ErrFriendRequestNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case friend.ErrCannotRejectNonPending:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "フレンドリクエストの拒否に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FriendHandler) CancelFriendRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
		return
	}

	// URLパスからリクエストIDを取得
	requestID := r.URL.Path[len("/api/friends/"):]
	if requestID == "" {
		respondError(w, http.StatusBadRequest, "リクエストIDが必要です")
		return
	}

	err := h.useCase.CancelFriendRequest(r.Context(), requestID, userID)
	if err != nil {
		switch err {
		case friend.ErrFriendRequestNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "フレンドリクエストのキャンセルに失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FriendHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
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

	friends, err := h.useCase.GetFriends(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "フレンド一覧の取得に失敗しました")
		return
	}

	responses := make([]FriendResponse, 0, len(friends))
	for _, f := range friends {
		responses = append(responses, toFriendResponse(f))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"friends": responses,
		"limit":   limit,
		"offset":  offset,
		"count":   len(responses),
	})
}

func (h *FriendHandler) GetPendingReceivedRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
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

	requests, err := h.useCase.GetPendingReceivedRequests(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "受信リクエスト一覧の取得に失敗しました")
		return
	}

	responses := make([]FriendResponse, 0, len(requests))
	for _, f := range requests {
		responses = append(responses, toFriendResponse(f))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"requests": responses,
		"limit":    limit,
		"offset":   offset,
		"count":    len(responses),
	})
}

func (h *FriendHandler) GetPendingSentRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
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

	requests, err := h.useCase.GetPendingSentRequests(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "送信リクエスト一覧の取得に失敗しました")
		return
	}

	responses := make([]FriendResponse, 0, len(requests))
	for _, f := range requests {
		responses = append(responses, toFriendResponse(f))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"requests": responses,
		"limit":    limit,
		"offset":   offset,
		"count":    len(responses),
	})
}

func (h *FriendHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// TODO: 実際の実装ではJWTトークンなどから現在のユーザーIDを取得する
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "ユーザー認証が必要です")
		return
	}

	friendUserID := r.URL.Query().Get("friend_user_id")
	if friendUserID == "" {
		respondError(w, http.StatusBadRequest, "フレンドのユーザーIDが必要です")
		return
	}

	err := h.useCase.RemoveFriend(r.Context(), userID, friendUserID)
	if err != nil {
		switch err {
		case friend.ErrFriendRequestNotFound:
			respondError(w, http.StatusNotFound, "フレンド関係が見つかりません")
		default:
			respondError(w, http.StatusInternalServerError, "フレンドの削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
