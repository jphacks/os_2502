package group

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// GroupType represents the type of group
type GroupType string

const (
	GroupTypeLocalTemporary  GroupType = "local_temporary"
	GroupTypeGlobalTemporary GroupType = "global_temporary"
	GroupTypePermanent       GroupType = "permanent"
)

// GroupStatus represents the status of group workflow
type GroupStatus string

const (
	GroupStatusRecruiting  GroupStatus = "recruiting"   // メンバー募集中
	GroupStatusReadyCheck  GroupStatus = "ready_check"  // メンバー確定、準備完了待ち
	GroupStatusCountdown   GroupStatus = "countdown"    // カウントダウン中
	GroupStatusPhotoTaking GroupStatus = "photo_taking" // 撮影中
	GroupStatusCompleted   GroupStatus = "completed"    // 完了
	GroupStatusExpired     GroupStatus = "expired"      // 期限切れ
)

const (
	// SystemMaxMember is the system-wide maximum number of members in a group
	SystemMaxMember = 100
)

type Group struct {
	id                    string
	ownerUserID           string
	name                  string
	groupType             GroupType
	status                GroupStatus
	maxMember             int
	currentMemberCount    int
	invitationToken       string
	finalizedAt           *time.Time
	countdownStartedAt    *time.Time
	scheduledCaptureTime  *time.Time
	templateID            *string
	expiresAt             *time.Time
	createdAt             time.Time
	updatedAt             time.Time
}

// NewGroup creates a new group (max_member is system fixed value)
func NewGroup(ownerUserID, name string, groupType GroupType, expiresAt *time.Time) (*Group, error) {
	if ownerUserID == "" {
		return nil, ErrInvalidOwnerUserID
	}

	// 文字数をルーン数でカウント
	nameLen := utf8.RuneCountInString(name)
	if nameLen == 0 || nameLen > 15 {
		return nil, ErrInvalidName
	}

	if !isValidGroupType(groupType) {
		return nil, ErrInvalidGroupType
	}

	now := time.Now()
	return &Group{
		id:                   uuid.New().String(),
		ownerUserID:          ownerUserID,
		name:                 name,
		groupType:            groupType,
		status:               GroupStatusRecruiting,
		maxMember:            SystemMaxMember, // システム固定値
		currentMemberCount:   0,
		invitationToken:      uuid.New().String(),
		finalizedAt:          nil,
		countdownStartedAt:   nil,
		scheduledCaptureTime: nil,
		templateID:           nil,
		expiresAt:            expiresAt,
		createdAt:            now,
		updatedAt:            now,
	}, nil
}

// Reconstruct は既存のグループを復元（リポジトリから取得時に使用）
func Reconstruct(
	id, ownerUserID, name string,
	groupType GroupType,
	status GroupStatus,
	maxMember, currentMemberCount int,
	invitationToken string,
	finalizedAt, countdownStartedAt, scheduledCaptureTime *time.Time,
	templateID *string,
	expiresAt *time.Time,
	createdAt, updatedAt time.Time,
) (*Group, error) {
	if id == "" {
		return nil, ErrInvalidGroupID
	}
	if ownerUserID == "" {
		return nil, ErrInvalidOwnerUserID
	}

	nameLen := utf8.RuneCountInString(name)
	if nameLen == 0 || nameLen > 15 {
		return nil, ErrInvalidName
	}

	if maxMember < 1 || maxMember > 100 {
		return nil, ErrInvalidMaxMember
	}

	if !isValidGroupType(groupType) {
		return nil, ErrInvalidGroupType
	}

	if !isValidGroupStatus(status) {
		return nil, ErrInvalidGroupStatus
	}

	return &Group{
		id:                   id,
		ownerUserID:          ownerUserID,
		name:                 name,
		groupType:            groupType,
		status:               status,
		maxMember:            maxMember,
		currentMemberCount:   currentMemberCount,
		invitationToken:      invitationToken,
		finalizedAt:          finalizedAt,
		countdownStartedAt:   countdownStartedAt,
		scheduledCaptureTime: scheduledCaptureTime,
		templateID:           templateID,
		expiresAt:            expiresAt,
		createdAt:            createdAt,
		updatedAt:            updatedAt,
	}, nil
}

// Getters
func (g *Group) ID() string {
	return g.id
}

func (g *Group) OwnerUserID() string {
	return g.ownerUserID
}

func (g *Group) Name() string {
	return g.name
}

func (g *Group) GroupType() GroupType {
	return g.groupType
}

func (g *Group) Status() GroupStatus {
	return g.status
}

func (g *Group) MaxMember() int {
	return g.maxMember
}

func (g *Group) CurrentMemberCount() int {
	return g.currentMemberCount
}

func (g *Group) InvitationToken() string {
	return g.invitationToken
}

func (g *Group) FinalizedAt() *time.Time {
	return g.finalizedAt
}

func (g *Group) CountdownStartedAt() *time.Time {
	return g.countdownStartedAt
}

func (g *Group) ScheduledCaptureTime() *time.Time {
	return g.scheduledCaptureTime
}

func (g *Group) TemplateID() *string {
	return g.templateID
}

func (g *Group) ExpiresAt() *time.Time {
	return g.expiresAt
}

func (g *Group) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Group) UpdatedAt() time.Time {
	return g.updatedAt
}

// Business logic methods

func (g *Group) UpdateName(name string) error {
	nameLen := utf8.RuneCountInString(name)
	if nameLen == 0 || nameLen > 15 {
		return ErrInvalidName
	}
	g.name = name
	g.updatedAt = time.Now()
	return nil
}

func (g *Group) UpdateMaxMember(maxMember int) error {
	if maxMember < 1 || maxMember > 100 {
		return ErrInvalidMaxMember
	}
	if maxMember < g.currentMemberCount {
		return ErrMaxMemberLessThanCurrent
	}
	g.maxMember = maxMember
	g.updatedAt = time.Now()
	return nil
}

// IncrementMemberCount increments the current member count
func (g *Group) IncrementMemberCount() error {
	if g.currentMemberCount >= g.maxMember {
		return ErrGroupFull
	}
	g.currentMemberCount++
	g.updatedAt = time.Now()
	return nil
}

// DecrementMemberCount decrements the current member count
func (g *Group) DecrementMemberCount() error {
	if g.currentMemberCount <= 0 {
		return ErrInvalidMemberCount
	}
	g.currentMemberCount--
	g.updatedAt = time.Now()
	return nil
}

// FinalizeMembers finalizes the group members and moves to ready_check status
// This sets the actual max_member to the current member count
func (g *Group) FinalizeMembers() error {
	if g.status != GroupStatusRecruiting {
		return ErrGroupNotRecruiting
	}
	if g.currentMemberCount == 0 {
		return ErrNoMembers
	}
	now := time.Now()
	g.status = GroupStatusReadyCheck
	g.maxMember = g.currentMemberCount // 実際のメンバー数に確定
	g.finalizedAt = &now
	g.updatedAt = now
	return nil
}

// StartCountdown starts the countdown and sets the scheduled capture time (10 seconds from now)
func (g *Group) StartCountdown(countdownSeconds int, templateID string) error {
	if g.status != GroupStatusReadyCheck {
		return ErrGroupNotReadyCheck
	}
	now := time.Now()
	scheduledTime := now.Add(time.Duration(countdownSeconds) * time.Second)

	g.status = GroupStatusCountdown
	g.countdownStartedAt = &now
	g.scheduledCaptureTime = &scheduledTime
	g.templateID = &templateID
	g.updatedAt = now
	return nil
}

// StartPhotoTaking moves to photo taking status
func (g *Group) StartPhotoTaking() error {
	if g.status != GroupStatusCountdown {
		return ErrGroupNotCountdown
	}
	g.status = GroupStatusPhotoTaking
	g.updatedAt = time.Now()
	return nil
}

// Complete completes the group
func (g *Group) Complete() error {
	if g.status != GroupStatusPhotoTaking {
		return ErrGroupNotPhotoTaking
	}
	g.status = GroupStatusCompleted
	g.updatedAt = time.Now()
	return nil
}

// Expire expires the group
func (g *Group) Expire() error {
	g.status = GroupStatusExpired
	g.updatedAt = time.Now()
	return nil
}

// IsExpired checks if the group is expired
func (g *Group) IsExpired() bool {
	if g.expiresAt == nil {
		return false
	}
	return time.Now().After(*g.expiresAt)
}

// IsFull checks if the group is full
func (g *Group) IsFull() bool {
	return g.currentMemberCount >= g.maxMember
}

// CanJoin checks if a new member can join
func (g *Group) CanJoin() bool {
	return g.status == GroupStatusRecruiting && !g.IsFull() && !g.IsExpired()
}

// Helper functions

func isValidGroupType(gt GroupType) bool {
	switch gt {
	case GroupTypeLocalTemporary, GroupTypeGlobalTemporary, GroupTypePermanent:
		return true
	default:
		return false
	}
}

func isValidGroupStatus(gs GroupStatus) bool {
	switch gs {
	case GroupStatusRecruiting, GroupStatusReadyCheck, GroupStatusCountdown,
		GroupStatusPhotoTaking, GroupStatusCompleted, GroupStatusExpired:
		return true
	default:
		return false
	}
}
