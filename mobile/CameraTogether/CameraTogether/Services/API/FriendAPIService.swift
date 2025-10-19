import Foundation

/// フレンド関連のAPI通信を管理するサービス
class FriendAPIService: APIServiceBase {
    static let shared = FriendAPIService()

    private override init() {
        super.init()
    }

    // MARK: - Friend APIs

    /// フレンドリクエスト送信
    /// - Parameters:
    ///   - requesterId: リクエスト送信者のユーザーID
    ///   - addresseeId: リクエスト受信者のユーザーID
    /// - Returns: 作成されたフレンドリクエスト
    func sendFriendRequest(requesterId: String, addresseeId: String) async throws -> APIFriend {
        let url = baseURL.appendingPathComponent("friends")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(requesterId, forHTTPHeaderField: "X-User-ID")

        let body = SendFriendRequestRequest(addresseeId: addresseeId)
        request.httpBody = try JSONEncoder().encode(body)

        return try await performRequest(request, expecting: APIFriend.self, successStatusCode: 201)
    }

    /// フレンド一覧取得
    /// - Parameters:
    ///   - userId: ユーザーID
    ///   - limit: 取得件数
    ///   - offset: オフセット
    /// - Returns: フレンド一覧
    func getFriends(userId: String, limit: Int = 100, offset: Int = 0) async throws
        -> [APIFriend]
    {
        var components = URLComponents(
            url: baseURL.appendingPathComponent("friends"), resolvingAgainstBaseURL: false)!
        components.queryItems = [
            URLQueryItem(name: "limit", value: "\(limit)"),
            URLQueryItem(name: "offset", value: "\(offset)"),
        ]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue(userId, forHTTPHeaderField: "X-User-ID")

        let response = try await performRequest(request, expecting: FriendListResponse.self)
        return response.friends
    }

    /// フレンドリクエスト承認
    /// - Parameters:
    ///   - requestId: リクエストID
    ///   - userId: 承認するユーザーID
    /// - Returns: 更新されたフレンド情報
    func acceptFriendRequest(requestId: String, userId: String) async throws -> APIFriend {
        let url = baseURL.appendingPathComponent("friends").appendingPathComponent(requestId)
            .appendingPathComponent("accept")
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue(userId, forHTTPHeaderField: "X-User-ID")

        return try await performRequest(request, expecting: APIFriend.self)
    }

    /// フレンドリクエスト拒否
    /// - Parameters:
    ///   - requestId: リクエストID
    ///   - userId: 拒否するユーザーID
    /// - Returns: 更新されたフレンド情報
    func rejectFriendRequest(requestId: String, userId: String) async throws -> APIFriend {
        let url = baseURL.appendingPathComponent("friends").appendingPathComponent(requestId)
            .appendingPathComponent("reject")
        var request = URLRequest(url: url)
        request.httpMethod = "PUT"
        request.setValue(userId, forHTTPHeaderField: "X-User-ID")

        return try await performRequest(request, expecting: APIFriend.self)
    }

    /// フレンド削除
    /// - Parameters:
    ///   - friendId: フレンドID
    ///   - userId: 削除するユーザーID
    func deleteFriend(friendId: String, userId: String) async throws {
        let url = baseURL.appendingPathComponent("friends").appendingPathComponent(friendId)
        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"
        request.setValue(userId, forHTTPHeaderField: "X-User-ID")

        try await performRequest(request, successStatusCode: 204)
    }
}
