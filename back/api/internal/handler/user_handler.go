package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/user"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type UserHandler struct {
	useCase *usecase.UserUseCase
}

func NewUserHandler(useCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

type CreateUserRequest struct {
	FirebaseUID string `json:"firebase_uid"`
	Name        string `json:"name"`
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	FirebaseUID string  `json:"firebase_uid"`
	Name        string  `json:"name"`
	Username    *string `json:"username,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// UserエンティティをUserResponseに変換
func toResponse(u *user.User) UserResponse {
	return UserResponse{
		ID:          u.ID().String(),
		FirebaseUID: u.FirebaseUID(),
		Name:        u.Name(),
		Username:    u.Username(),
		CreatedAt:   u.CreatedAt().Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   u.UpdatedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	u, err := h.useCase.CreateUser(r.Context(), req.FirebaseUID, req.Name)
	if err != nil {
		switch err {
		case user.ErrInvalidFirebaseUID, user.ErrInvalidName:
			respondError(w, http.StatusBadRequest, err.Error())
		case user.ErrUserAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "ユーザーの作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toResponse(u))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/users/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	u, err := h.useCase.GetUserByID(r.Context(), id)
	if err != nil {
		if err == user.ErrUserNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "ユーザーの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) GetUserByFirebaseUID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	firebaseUID := r.URL.Query().Get("firebase_uid")
	if firebaseUID == "" {
		respondError(w, http.StatusBadRequest, "Firebase UIDが必要です")
		return
	}

	u, err := h.useCase.GetUserByFirebaseUID(r.Context(), firebaseUID)
	if err != nil {
		if err == user.ErrUserNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "ユーザーの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/users/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	u, err := h.useCase.UpdateUserName(r.Context(), id, req.Name)
	if err != nil {
		switch err {
		case user.ErrInvalidName:
			respondError(w, http.StatusBadRequest, err.Error())
		case user.ErrUserNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "ユーザーの更新に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/users/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	if err := h.useCase.DeleteUser(r.Context(), id); err != nil {
		if err == user.ErrUserNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "ユーザーの削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
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

	users, err := h.useCase.ListUsers(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ユーザー一覧の取得に失敗しました")
		return
	}

	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, toResponse(u))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"users":  responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

func (h *UserHandler) SetUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/users/"):]
	if idStr == "" || len(idStr) < len("/username")+1 {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}
	// "/username" を削除
	idStr = idStr[:len(idStr)-len("/username")]

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	u, err := h.useCase.SetUsername(r.Context(), id, req.Username)
	if err != nil {
		switch err {
		case user.ErrInvalidUsername:
			respondError(w, http.StatusBadRequest, err.Error())
		case user.ErrUsernameAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		case user.ErrUserNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "ユーザー名の設定に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		respondError(w, http.StatusBadRequest, "ユーザー名が必要です")
		return
	}

	u, err := h.useCase.GetUserByUsername(r.Context(), username)
	if err != nil {
		if err == user.ErrUserNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "ユーザーの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toResponse(u))
}

func (h *UserHandler) SearchUsersByUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		respondError(w, http.StatusBadRequest, "検索クエリが必要です")
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

	users, err := h.useCase.SearchUsersByUsername(r.Context(), query, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "ユーザー検索に失敗しました")
		return
	}

	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, toResponse(u))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"users":  responses,
		"limit":  limit,
		"offset": offset,
		"count":  len(responses),
	})
}

// respondJSON JSONレスポンスを返す
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// respondError エラーレスポンスを返す
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
