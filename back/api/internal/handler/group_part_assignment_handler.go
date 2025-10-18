package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_part_assignment"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type GroupPartAssignmentHandler struct {
	useCase *usecase.GroupPartAssignmentUseCase
}

func NewGroupPartAssignmentHandler(useCase *usecase.GroupPartAssignmentUseCase) *GroupPartAssignmentHandler {
	return &GroupPartAssignmentHandler{useCase: useCase}
}

type CreateGroupPartAssignmentRequest struct {
	GroupID    string `json:"group_id"`
	UserID     string `json:"user_id"`
	PartID     string `json:"part_id"`
	CollageDay string `json:"collage_day"`
}

type GroupPartAssignmentResponse struct {
	AssignmentID string `json:"assignment_id"`
	GroupID      string `json:"group_id"`
	UserID       string `json:"user_id"`
	PartID       string `json:"part_id"`
	CollageDay   string `json:"collage_day"`
	AssignedAt   string `json:"assigned_at"`
}

// GroupPartAssignmentエンティティをGroupPartAssignmentResponseに変換
func toGroupPartAssignmentResponse(gpa *group_part_assignment.GroupPartAssignment) GroupPartAssignmentResponse {
	return GroupPartAssignmentResponse{
		AssignmentID: gpa.AssignmentID().String(),
		GroupID:      gpa.GroupID(),
		UserID:       gpa.UserID().String(),
		PartID:       gpa.PartID().String(),
		CollageDay:   gpa.CollageDay().Format("2006-01-02"),
		AssignedAt:   gpa.AssignedAt().Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (h *GroupPartAssignmentHandler) CreateGroupPartAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	var req CreateGroupPartAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	partID, err := uuid.Parse(req.PartID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なパーツIDです")
		return
	}

	collageDay, err := time.Parse("2006-01-02", req.CollageDay)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なコラージュ日です (形式: YYYY-MM-DD)")
		return
	}

	gpa, err := h.useCase.CreateGroupPartAssignment(r.Context(), req.GroupID, userID, partID, collageDay)
	if err != nil {
		switch err {
		case group_part_assignment.ErrInvalidGroupID, group_part_assignment.ErrInvalidUserID, group_part_assignment.ErrInvalidPartID, group_part_assignment.ErrInvalidCollageDay:
			respondError(w, http.StatusBadRequest, err.Error())
		case group_part_assignment.ErrGroupPartAssignmentAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "グループパーツ割り当ての作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toGroupPartAssignmentResponse(gpa))
}

func (h *GroupPartAssignmentHandler) GetGroupPartAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/group-part-assignments/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "グループパーツ割り当てIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なグループパーツ割り当てIDです")
		return
	}

	gpa, err := h.useCase.GetGroupPartAssignmentByID(r.Context(), id)
	if err != nil {
		if err == group_part_assignment.ErrGroupPartAssignmentNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "グループパーツ割り当ての取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toGroupPartAssignmentResponse(gpa))
}

func (h *GroupPartAssignmentHandler) DeleteGroupPartAssignment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	// URLパスからIDを取得
	idStr := r.URL.Path[len("/api/group-part-assignments/"):]
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "グループパーツ割り当てIDが必要です")
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なグループパーツ割り当てIDです")
		return
	}

	if err := h.useCase.DeleteGroupPartAssignment(r.Context(), id); err != nil {
		if err == group_part_assignment.ErrGroupPartAssignmentNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "グループパーツ割り当ての削除に失敗しました")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupPartAssignmentHandler) ListGroupPartAssignments(w http.ResponseWriter, r *http.Request) {
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

	assignments, err := h.useCase.ListGroupPartAssignments(r.Context(), limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "グループパーツ割り当て一覧の取得に失敗しました")
		return
	}

	responses := make([]GroupPartAssignmentResponse, 0, len(assignments))
	for _, gpa := range assignments {
		responses = append(responses, toGroupPartAssignmentResponse(gpa))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"group_part_assignments": responses,
		"limit":                  limit,
		"offset":                 offset,
		"count":                  len(responses),
	})
}

func (h *GroupPartAssignmentHandler) GetGroupPartAssignmentsByGroupAndDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	collageDayStr := r.URL.Query().Get("collage_day")
	if collageDayStr == "" {
		respondError(w, http.StatusBadRequest, "コラージュ日が必要です")
		return
	}

	collageDay, err := time.Parse("2006-01-02", collageDayStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なコラージュ日です (形式: YYYY-MM-DD)")
		return
	}

	assignments, err := h.useCase.GetGroupPartAssignmentsByGroupAndDay(r.Context(), groupID, collageDay)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "グループパーツ割り当ての取得に失敗しました")
		return
	}

	responses := make([]GroupPartAssignmentResponse, 0, len(assignments))
	for _, gpa := range assignments {
		responses = append(responses, toGroupPartAssignmentResponse(gpa))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"group_part_assignments": responses,
		"count":                  len(responses),
	})
}

func (h *GroupPartAssignmentHandler) GetGroupPartAssignmentByUserGroupAndDay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なユーザーIDです")
		return
	}

	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	collageDayStr := r.URL.Query().Get("collage_day")
	if collageDayStr == "" {
		respondError(w, http.StatusBadRequest, "コラージュ日が必要です")
		return
	}

	collageDay, err := time.Parse("2006-01-02", collageDayStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なコラージュ日です (形式: YYYY-MM-DD)")
		return
	}

	gpa, err := h.useCase.GetGroupPartAssignmentByUserGroupAndDay(r.Context(), userID, groupID, collageDay)
	if err != nil {
		if err == group_part_assignment.ErrGroupPartAssignmentNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "グループパーツ割り当ての取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toGroupPartAssignmentResponse(gpa))
}

func (h *GroupPartAssignmentHandler) GetGroupPartAssignmentsByPartID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "メソッドが許可されていません")
		return
	}

	partIDStr := r.URL.Query().Get("part_id")
	if partIDStr == "" {
		respondError(w, http.StatusBadRequest, "パーツIDが必要です")
		return
	}

	partID, err := uuid.Parse(partIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "無効なパーツIDです")
		return
	}

	assignments, err := h.useCase.GetGroupPartAssignmentsByPartID(r.Context(), partID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "グループパーツ割り当ての取得に失敗しました")
		return
	}

	responses := make([]GroupPartAssignmentResponse, 0, len(assignments))
	for _, gpa := range assignments {
		responses = append(responses, toGroupPartAssignmentResponse(gpa))
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"group_part_assignments": responses,
		"count":                  len(responses),
	})
}
