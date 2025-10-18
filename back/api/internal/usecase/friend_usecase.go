package usecase

import (
	"context"

	"github.com/jphacks/os_2502/back/api/internal/domain/friend"
)

type FriendUseCase struct {
	repo friend.Repository
}

func NewFriendUseCase(repo friend.Repository) *FriendUseCase {
	return &FriendUseCase{repo: repo}
}

// SendFriendRequest sends a friend request
func (uc *FriendUseCase) SendFriendRequest(ctx context.Context, requesterID, addresseeID string) (*friend.Friend, error) {
	// 既存のフレンドリクエストまたはフレンド関係をチェック
	existing, err := uc.repo.FindByRequesterAndAddressee(ctx, requesterID, addresseeID)
	if err == nil && existing != nil {
		if existing.IsAccepted() {
			return nil, friend.ErrAlreadyFriends
		}
		if existing.IsPending() {
			return nil, friend.ErrFriendRequestAlreadyExists
		}
	}

	// 逆方向のリクエストもチェック
	reverse, err := uc.repo.FindByRequesterAndAddressee(ctx, addresseeID, requesterID)
	if err == nil && reverse != nil {
		if reverse.IsAccepted() {
			return nil, friend.ErrAlreadyFriends
		}
		if reverse.IsPending() {
			return nil, friend.ErrFriendRequestAlreadyExists
		}
	}

	// 新しいフレンドリクエストを作成
	f, err := friend.NewFriend(requesterID, addresseeID)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

// AcceptFriendRequest accepts a friend request
func (uc *FriendUseCase) AcceptFriendRequest(ctx context.Context, requestID, userID string) (*friend.Friend, error) {
	// フレンドリクエストを取得
	f, err := uc.repo.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, friend.ErrFriendRequestNotFound
	}

	// 受信者本人かチェック
	if f.AddresseeID() != userID {
		return nil, friend.ErrFriendRequestNotFound
	}

	// 承認
	if err := f.Accept(); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(ctx, f); err != nil {
		return nil, err
	}

	return f, nil
}

// RejectFriendRequest rejects a friend request
func (uc *FriendUseCase) RejectFriendRequest(ctx context.Context, requestID, userID string) error {
	// フレンドリクエストを取得
	f, err := uc.repo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if f == nil {
		return friend.ErrFriendRequestNotFound
	}

	// 受信者本人かチェック
	if f.AddresseeID() != userID {
		return friend.ErrFriendRequestNotFound
	}

	// 拒否
	if err := f.Reject(); err != nil {
		return err
	}

	// 拒否されたリクエストは削除
	return uc.repo.Delete(ctx, requestID)
}

// CancelFriendRequest cancels a friend request
func (uc *FriendUseCase) CancelFriendRequest(ctx context.Context, requestID, userID string) error {
	// フレンドリクエストを取得
	f, err := uc.repo.FindByID(ctx, requestID)
	if err != nil {
		return err
	}
	if f == nil {
		return friend.ErrFriendRequestNotFound
	}

	// 送信者本人かチェック
	if f.RequesterID() != userID {
		return friend.ErrFriendRequestNotFound
	}

	// キャンセル（削除）
	return uc.repo.Delete(ctx, requestID)
}

// GetFriends retrieves all friends for a user
func (uc *FriendUseCase) GetFriends(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindAcceptedFriends(ctx, userID, limit, offset)
}

// GetPendingReceivedRequests retrieves pending received friend requests
func (uc *FriendUseCase) GetPendingReceivedRequests(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindPendingReceivedRequests(ctx, userID, limit, offset)
}

// GetPendingSentRequests retrieves pending sent friend requests
func (uc *FriendUseCase) GetPendingSentRequests(ctx context.Context, userID string, limit, offset int) ([]*friend.Friend, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return uc.repo.FindPendingSentRequests(ctx, userID, limit, offset)
}

// RemoveFriend removes a friend
func (uc *FriendUseCase) RemoveFriend(ctx context.Context, userID, friendUserID string) error {
	// フレンド関係を検索（双方向）
	f, err := uc.repo.FindByRequesterAndAddressee(ctx, userID, friendUserID)
	if err != nil || f == nil {
		// 逆方向も確認
		f, err = uc.repo.FindByRequesterAndAddressee(ctx, friendUserID, userID)
		if err != nil || f == nil {
			return friend.ErrFriendRequestNotFound
		}
	}

	// フレンド関係を削除
	return uc.repo.Delete(ctx, f.ID())
}
