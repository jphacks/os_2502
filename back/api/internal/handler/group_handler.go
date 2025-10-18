package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jphacks/os_2502/back/api/internal/domain/group"
	"github.com/jphacks/os_2502/back/api/internal/domain/group_member"
	"github.com/jphacks/os_2502/back/api/internal/usecase"
)

type GroupHandler struct {
	useCase *usecase.GroupUseCase
}

func NewGroupHandler(useCase *usecase.GroupUseCase) *GroupHandler {
	return &GroupHandler{useCase: useCase}
}

// Request/Response types

type CreateGroupRequest struct {
	OwnerUserID string `json:"owner_user_id"`
	Name        string `json:"name"`
	GroupType   string `json:"group_type"`           // "local_temporary", "global_temporary", "permanent"
	ExpiresAt   string `json:"expires_at,omitempty"` // ISO 8601 format
}

type JoinGroupRequest struct {
	UserID string `json:"user_id"`
}

type FinalizeGroupRequest struct {
	UserID string `json:"user_id"`
}

type MarkReadyRequest struct {
	UserID string `json:"user_id"`
}

type GroupResponse struct {
	ID                 string  `json:"id"`
	OwnerUserID        string  `json:"owner_user_id"`
	Name               string  `json:"name"`
	GroupType          string  `json:"group_type"`
	Status             string  `json:"status"`
	MaxMember          int     `json:"max_member"`
	CurrentMemberCount int     `json:"current_member_count"`
	InvitationToken    string  `json:"invitation_token"`
	FinalizedAt        *string `json:"finalized_at,omitempty"`
	CountdownStartedAt *string `json:"countdown_started_at,omitempty"`
	ExpiresAt          *string `json:"expires_at,omitempty"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}

type GroupMemberResponse struct {
	ID          string  `json:"id"`
	GroupID     string  `json:"group_id"`
	UserID      string  `json:"user_id"`
	IsOwner     bool    `json:"is_owner"`
	ReadyStatus bool    `json:"ready_status"`
	ReadyAt     *string `json:"ready_at,omitempty"`
	JoinedAt    string  `json:"joined_at"`
}

type GroupListResponse struct {
	Groups     []GroupResponse `json:"groups"`
	TotalCount int             `json:"total_count"`
}

// Conversion functions

func toGroupResponse(g *group.Group) GroupResponse {
	resp := GroupResponse{
		ID:                 g.ID(),
		OwnerUserID:        g.OwnerUserID(),
		Name:               g.Name(),
		GroupType:          string(g.GroupType()),
		Status:             string(g.Status()),
		MaxMember:          g.MaxMember(),
		CurrentMemberCount: g.CurrentMemberCount(),
		InvitationToken:    g.InvitationToken(),
		CreatedAt:          g.CreatedAt().Format(time.RFC3339),
		UpdatedAt:          g.UpdatedAt().Format(time.RFC3339),
	}

	if finalizedAt := g.FinalizedAt(); finalizedAt != nil {
		str := finalizedAt.Format(time.RFC3339)
		resp.FinalizedAt = &str
	}

	if countdownStartedAt := g.CountdownStartedAt(); countdownStartedAt != nil {
		str := countdownStartedAt.Format(time.RFC3339)
		resp.CountdownStartedAt = &str
	}

	if expiresAt := g.ExpiresAt(); expiresAt != nil {
		str := expiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &str
	}

	return resp
}

func toGroupMemberResponse(m *group_member.GroupMember) GroupMemberResponse {
	resp := GroupMemberResponse{
		ID:          m.ID(),
		GroupID:     m.GroupID(),
		UserID:      m.UserID(),
		IsOwner:     m.IsOwner(),
		ReadyStatus: m.ReadyStatus(),
		JoinedAt:    m.JoinedAt().Format(time.RFC3339),
	}

	if readyAt := m.ReadyAt(); readyAt != nil {
		str := readyAt.Format(time.RFC3339)
		resp.ReadyAt = &str
	}

	return resp
}

// Handlers

// CreateGroup creates a new group
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req CreateGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	// Parse group type
	var groupType group.GroupType
	switch req.GroupType {
	case "local_temporary":
		groupType = group.GroupTypeLocalTemporary
	case "global_temporary":
		groupType = group.GroupTypeGlobalTemporary
	case "permanent":
		groupType = group.GroupTypePermanent
	case "":
		groupType = group.GroupTypeGlobalTemporary // default
	default:
		respondError(w, http.StatusBadRequest, "無効なグループタイプです")
		return
	}

	// Parse expires_at
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			respondError(w, http.StatusBadRequest, "無効な有効期限です")
			return
		}
		expiresAt = &t
	}

	g, err := h.useCase.CreateGroup(r.Context(), req.OwnerUserID, req.Name, groupType, expiresAt)
	if err != nil {
		switch err {
		case group.ErrInvalidOwnerUserID, group.ErrInvalidName, group.ErrInvalidGroupType:
			respondError(w, http.StatusBadRequest, err.Error())
		case group.ErrGroupAlreadyExists:
			respondError(w, http.StatusConflict, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "グループの作成に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusCreated, toGroupResponse(g))
}

// GetGroupByID retrieves a group by ID
func (h *GroupHandler) GetGroupByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	if id == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	g, err := h.useCase.GetGroupByID(r.Context(), id)
	if err != nil {
		if err == group.ErrGroupNotFound {
			respondError(w, http.StatusNotFound, err.Error())
		} else {
			respondError(w, http.StatusInternalServerError, "グループの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toGroupResponse(g))
}

// GetGroupByInvitationToken retrieves a group by invitation token
func (h *GroupHandler) GetGroupByInvitationToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("invitation_token")
	if token == "" {
		respondError(w, http.StatusBadRequest, "招待トークンが必要です")
		return
	}

	g, err := h.useCase.GetGroupByInvitationToken(r.Context(), token)
	if err != nil {
		if err == group.ErrGroupNotFound {
			respondError(w, http.StatusNotFound, "グループが見つかりません")
		} else {
			respondError(w, http.StatusInternalServerError, "グループの取得に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toGroupResponse(g))
}

// GetGroupsByOwnerUserID retrieves groups by owner user ID
func (h *GroupHandler) GetGroupsByOwnerUserID(w http.ResponseWriter, r *http.Request) {
	ownerUserID := r.URL.Query().Get("owner_user_id")
	if ownerUserID == "" {
		respondError(w, http.StatusBadRequest, "オーナーユーザーIDが必要です")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	groups, err := h.useCase.GetGroupsByOwnerUserID(r.Context(), ownerUserID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "グループの取得に失敗しました")
		return
	}

	groupResponses := make([]GroupResponse, len(groups))
	for i, g := range groups {
		groupResponses[i] = toGroupResponse(g)
	}

	respondJSON(w, http.StatusOK, GroupListResponse{
		Groups:     groupResponses,
		TotalCount: len(groupResponses),
	})
}

// JoinGroup joins a group via invitation token
func (h *GroupHandler) JoinGroup(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/api/groups/join/")
	if token == "" {
		respondError(w, http.StatusBadRequest, "招待トークンが必要です")
		return
	}

	var req JoinGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	g, err := h.useCase.JoinGroup(r.Context(), token, req.UserID)
	if err != nil {
		switch err {
		case group.ErrGroupNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case group.ErrGroupFull:
			respondError(w, http.StatusBadRequest, err.Error())
		case group.ErrGroupExpired:
			respondError(w, http.StatusBadRequest, err.Error())
		case group.ErrGroupNotRecruiting:
			respondError(w, http.StatusBadRequest, err.Error())
		case group_member.ErrMemberAlreadyExists:
			respondError(w, http.StatusConflict, "既にグループに参加しています")
		default:
			respondError(w, http.StatusInternalServerError, "グループへの参加に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toGroupResponse(g))
}

// FinalizeGroupMembers finalizes group members (owner only)
func (h *GroupHandler) FinalizeGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/groups/"), "/finalize")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	var req FinalizeGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	g, err := h.useCase.FinalizeGroupMembers(r.Context(), groupID, req.UserID)
	if err != nil {
		switch err {
		case group.ErrGroupNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case group.ErrInvalidOwnerUserID:
			respondError(w, http.StatusForbidden, "オーナーのみ実行できます")
		case group.ErrGroupNotRecruiting:
			respondError(w, http.StatusBadRequest, err.Error())
		case group.ErrNoMembers:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "メンバー確定に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, toGroupResponse(g))
}

// MarkMemberReady marks a member as ready
func (h *GroupHandler) MarkMemberReady(w http.ResponseWriter, r *http.Request) {
	groupID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/groups/"), "/ready")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	var req MarkReadyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "リクエストボディが無効です")
		return
	}

	err := h.useCase.MarkMemberReady(r.Context(), groupID, req.UserID)
	if err != nil {
		switch err {
		case group_member.ErrMemberNotFound:
			respondError(w, http.StatusNotFound, "メンバーが見つかりません")
		case group_member.ErrAlreadyReady:
			respondError(w, http.StatusBadRequest, err.Error())
		default:
			respondError(w, http.StatusInternalServerError, "準備完了の設定に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "準備完了にしました"})
}

// GetGroupMembers retrieves all members of a group
func (h *GroupHandler) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	groupID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/groups/"), "/members")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	members, err := h.useCase.GetGroupMembers(r.Context(), groupID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "メンバーの取得に失敗しました")
		return
	}

	memberResponses := make([]GroupMemberResponse, len(members))
	for i, m := range members {
		memberResponses[i] = toGroupMemberResponse(m)
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"members": memberResponses,
		"count":   len(memberResponses),
	})
}

// LeaveGroup allows a member to leave a group
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	groupID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/groups/"), "/leave")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}

	err := h.useCase.LeaveGroup(r.Context(), groupID, userID)
	if err != nil {
		switch err {
		case group.ErrGroupNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case group.ErrInvalidOwnerUserID:
			respondError(w, http.StatusForbidden, "オーナーはグループを離脱できません")
		case group.ErrGroupNotRecruiting:
			respondError(w, http.StatusBadRequest, "メンバー募集中以外は離脱できません")
		case group_member.ErrMemberNotFound:
			respondError(w, http.StatusNotFound, "メンバーが見つかりません")
		default:
			respondError(w, http.StatusInternalServerError, "グループ離脱に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "グループを離脱しました"})
}

// DeleteGroup deletes a group (owner only)
func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupID := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	if groupID == "" {
		respondError(w, http.StatusBadRequest, "グループIDが必要です")
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "ユーザーIDが必要です")
		return
	}

	err := h.useCase.DeleteGroup(r.Context(), groupID, userID)
	if err != nil {
		switch err {
		case group.ErrGroupNotFound:
			respondError(w, http.StatusNotFound, err.Error())
		case group.ErrInvalidOwnerUserID:
			respondError(w, http.StatusForbidden, "オーナーのみ削除できます")
		default:
			respondError(w, http.StatusInternalServerError, "グループの削除に失敗しました")
		}
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "グループを削除しました"})
}

// ListGroups retrieves all groups, optionally filtered by owner_user_id
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	ownerUserID := r.URL.Query().Get("owner_user_id")

	var groups []*group.Group
	var err error

	// owner_user_idが指定されている場合はそのユーザーのグループのみを取得
	if ownerUserID != "" {
		groups, err = h.useCase.GetGroupsByOwnerUserID(r.Context(), ownerUserID, limit, offset)
	} else {
		groups, err = h.useCase.ListGroups(r.Context(), limit, offset)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, "グループの取得に失敗しました")
		return
	}

	groupResponses := make([]GroupResponse, len(groups))
	for i, g := range groups {
		groupResponses[i] = toGroupResponse(g)
	}

	respondJSON(w, http.StatusOK, GroupListResponse{
		Groups:     groupResponses,
		TotalCount: len(groupResponses),
	})
}
