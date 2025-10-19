import Foundation

@Observable
class FriendListViewModel {
    private var friends: [Friend] = []
    private var currentUserId: String = ""
    private let friendAPI = FriendAPIService.shared
    var isLoading: Bool = false
    var errorMessage: String?

    init() {
        self.friends = []
    }

    func setCurrentUserId(_ userId: String) {
        self.currentUserId = userId
    }

    func getFriends() -> [Friend] {
        friends.filter { $0.status == .accepted }
    }

    func getPendingRequests() -> [Friend] {
        friends.filter { $0.status == .pending }
    }

    /// フレンド一覧を取得
    @MainActor
    func fetchFriends() async {
        guard !currentUserId.isEmpty else { return }

        isLoading = true
        errorMessage = nil

        do {
            let apiFriends = try await friendAPI.getFriends(userId: currentUserId)

            // APIFriendをFriendに変換
            let convertedFriends = apiFriends.compactMap { apiFriend in
                Friend(from: apiFriend, currentUserId: currentUserId)
            }

            friends = convertedFriends
            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "フレンド一覧の取得に失敗しました: \(error.localizedDescription)"
            print("Error fetching friends: \(error)")
        }
    }

    /// フレンドリクエストを承認
    @MainActor
    func acceptFriend(_ friend: Friend) async {
        isLoading = true
        errorMessage = nil

        do {
            let updatedFriend = try await friendAPI.acceptFriendRequest(
                requestId: friend.id,
                userId: currentUserId
            )

            // ローカルのリストを更新
            if let index = friends.firstIndex(where: { $0.id == friend.id }),
                let converted = Friend(from: updatedFriend, currentUserId: currentUserId)
            {
                friends[index] = converted
            }

            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "フレンドリクエストの承認に失敗しました: \(error.localizedDescription)"
            print("Error accepting friend: \(error)")
        }
    }

    /// フレンドリクエストを拒否
    @MainActor
    func rejectFriend(_ friend: Friend) async {
        isLoading = true
        errorMessage = nil

        do {
            _ = try await friendAPI.rejectFriendRequest(
                requestId: friend.id,
                userId: currentUserId
            )

            // ローカルのリストから削除
            friends.removeAll { $0.id == friend.id }

            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "フレンドリクエストの拒否に失敗しました: \(error.localizedDescription)"
            print("Error rejecting friend: \(error)")
        }
    }

    /// フレンドリクエストを送信
    @MainActor
    func sendFriendRequest(to addresseeId: String) async {
        guard !currentUserId.isEmpty else { return }

        isLoading = true
        errorMessage = nil

        do {
            let newFriend = try await friendAPI.sendFriendRequest(
                requesterId: currentUserId,
                addresseeId: addresseeId
            )

            // ローカルのリストに追加
            if let converted = Friend(from: newFriend, currentUserId: currentUserId) {
                friends.append(converted)
            }

            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "フレンドリクエストの送信に失敗しました: \(error.localizedDescription)"
            print("Error sending friend request: \(error)")
        }
    }

    /// フレンドを削除
    @MainActor
    func deleteFriend(_ friend: Friend) async {
        isLoading = true
        errorMessage = nil

        do {
            try await friendAPI.deleteFriend(friendId: friend.id, userId: currentUserId)

            // ローカルのリストから削除
            friends.removeAll { $0.id == friend.id }

            isLoading = false
        } catch {
            isLoading = false
            errorMessage = "フレンドの削除に失敗しました: \(error.localizedDescription)"
            print("Error deleting friend: \(error)")
        }
    }
}
