import Foundation
import UIKit

/// グループ関連のAPI通信を管理するサービス
class GroupAPIService: APIServiceBase {
    static let shared = GroupAPIService()

    private override init() {
        super.init()
    }

    // MARK: - Group APIs

    /// グループ一覧取得
    /// - Parameters:
    ///   - ownerUserId: オーナーのユーザーID（オプション）
    ///   - limit: 取得件数（デフォルト10）
    ///   - offset: オフセット（デフォルト0）
    /// - Returns: グループの配列
    func getGroups(ownerUserId: String? = nil, limit: Int = 10, offset: Int = 0) async throws
        -> [APIGroup]
    {
        // モックデータを使用する場合
        if AppConfig.useMockData {
            let response = try MockDataService.shared.getGroups()
            return response.groups
        }

        // 実際のAPI呼び出し
        var components = URLComponents(
            url: baseURL.appendingPathComponent("groups"), resolvingAgainstBaseURL: false)!
        var queryItems: [URLQueryItem] = [
            URLQueryItem(name: "limit", value: "\(limit)"),
            URLQueryItem(name: "offset", value: "\(offset)"),
        ]
        if let ownerUserId = ownerUserId {
            queryItems.append(URLQueryItem(name: "owner_user_id", value: ownerUserId))
        }
        components.queryItems = queryItems

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        let groupList = try await performRequest(request, expecting: GroupListResponse.self)
        return groupList.groups
    }

    /// グループ作成
    /// - Parameters:
    ///   - ownerUserId: オーナーのユーザーID
    ///   - name: グループ名
    ///   - groupType: グループタイプ（デフォルト: global_temporary）
    /// - Returns: 作成されたグループ
    func createGroup(ownerUserId: String, name: String, groupType: String = "global_temporary")
        async throws -> APIGroup
    {
        let url = baseURL.appendingPathComponent("groups")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "owner_user_id": ownerUserId,
            "name": name,
            "group_type": groupType,
        ]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        return try await performRequest(request, expecting: APIGroup.self, successStatusCode: 201)
    }

    /// グループ取得
    /// - Parameter id: グループID
    /// - Returns: グループ情報
    func getGroup(id: String) async throws -> APIGroup {
        let url = baseURL.appendingPathComponent("groups").appendingPathComponent(id)
        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        return try await performRequest(request, expecting: APIGroup.self)
    }

    /// 招待トークンでグループ取得
    /// - Parameter invitationToken: 招待トークン
    /// - Returns: グループ情報
    func getGroupByInvitationToken(invitationToken: String) async throws -> APIGroup {
        var components = URLComponents(
            url: baseURL.appendingPathComponent("groups/by-invitation"),
            resolvingAgainstBaseURL: false)!
        components.queryItems = [URLQueryItem(name: "invitation_token", value: invitationToken)]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        return try await performRequest(request, expecting: APIGroup.self)
    }

    /// グループに参加
    /// - Parameters:
    ///   - token: 招待トークン
    ///   - userId: 参加するユーザーID
    /// - Returns: 参加後のグループ情報
    func joinGroup(token: String, userId: String) async throws -> APIGroup {
        let url = baseURL.appendingPathComponent("groups/join").appendingPathComponent(token)
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = JoinGroupRequest(userId: userId)
        request.httpBody = try JSONEncoder().encode(body)

        return try await performRequest(request, expecting: APIGroup.self)
    }

    /// グループメンバー一覧取得
    /// - Parameter groupId: グループID
    /// - Returns: メンバー一覧
    func getGroupMembers(groupId: String) async throws -> [GroupMember] {
        let url = baseURL.appendingPathComponent("groups").appendingPathComponent(groupId)
            .appendingPathComponent("members")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        let response = try await performRequest(request, expecting: GroupMemberListResponse.self)
        return response.members
    }

    /// メンバー確定
    /// - Parameters:
    ///   - groupId: グループID
    ///   - userId: オーナーのユーザーID
    /// - Returns: 更新後のグループ情報
    func finalizeGroupMembers(groupId: String, userId: String) async throws -> APIGroup {
        let url = baseURL.appendingPathComponent("groups").appendingPathComponent(groupId)
            .appendingPathComponent("finalize")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = FinalizeGroupRequest(userId: userId)
        request.httpBody = try JSONEncoder().encode(body)

        return try await performRequest(request, expecting: APIGroup.self)
    }

    /// カウントダウン開始
    /// - Parameters:
    ///   - groupId: グループID
    ///   - userId: オーナーのユーザーID
    /// - Returns: 更新されたグループ（撮影時刻を含む）
    func startCountdown(groupId: String, userId: String) async throws -> APIGroup {
        let url = baseURL.appendingPathComponent("groups").appendingPathComponent(groupId)
            .appendingPathComponent("start-countdown")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = ["user_id": userId]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        return try await performRequest(request, expecting: APIGroup.self)
    }

    /// 準備完了
    /// - Parameters:
    ///   - groupId: グループID
    ///   - userId: ユーザーID
    func markMemberReady(groupId: String, userId: String) async throws {
        let url = baseURL.appendingPathComponent("groups").appendingPathComponent(groupId)
            .appendingPathComponent("ready")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body = MarkReadyRequest(userId: userId)
        request.httpBody = try JSONEncoder().encode(body)

        try await performRequest(request, successStatusCode: 200)
    }

    /// グループ離脱
    /// - Parameters:
    ///   - groupId: グループID
    ///   - userId: ユーザーID
    func leaveGroup(groupId: String, userId: String) async throws {
        var components = URLComponents(
            url: baseURL.appendingPathComponent("groups").appendingPathComponent(groupId)
                .appendingPathComponent("leave"), resolvingAgainstBaseURL: false)!
        components.queryItems = [URLQueryItem(name: "user_id", value: userId)]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"

        try await performRequest(request, successStatusCode: 200)
    }

    /// グループ削除
    /// - Parameters:
    ///   - id: グループID
    ///   - userId: ユーザーID（オーナー確認用）
    func deleteGroup(id: String, userId: String) async throws {
        var components = URLComponents(
            url: baseURL.appendingPathComponent("groups").appendingPathComponent(id),
            resolvingAgainstBaseURL: false)!
        components.queryItems = [URLQueryItem(name: "user_id", value: userId)]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"

        try await performRequest(request, successStatusCode: 200)
    }

    // MARK: - Photo Upload

    /// 撮影した写真をサーバーにアップロード
    /// - Parameters:
    ///   - groupId: グループID
    ///   - userId: ユーザーID
    ///   - image: アップロードする画像
    ///   - frameIndex: フレームインデックス（担当パート番号）
    /// - Returns: アップロード結果
    func uploadPhoto(groupId: String, userId: String, image: UIImage, frameIndex: Int) async throws {

        // 画像をJPEGに変換
        guard let imageData = image.jpegData(compressionQuality: 0.8) else {
            throw APIError.invalidResponse
        }

        // 画像アップロードは /image エンドポイントを使用
        let imageBaseURL = URL(string: Configuration.shared.apiUrl)!.appendingPathComponent("image")
        let url = imageBaseURL
            .appendingPathComponent("groups")
            .appendingPathComponent(groupId)
            .appendingPathComponent("photos")

        // マルチパートフォームデータを作成
        let boundary = UUID().uuidString
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")

        var body = Data()

        // user_id フィールド
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"user_id\"\r\n\r\n".data(using: .utf8)!)
        body.append("\(userId)\r\n".data(using: .utf8)!)

        // frame_index フィールド
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"frame_index\"\r\n\r\n".data(using: .utf8)!)
        body.append("\(frameIndex)\r\n".data(using: .utf8)!)

        // photo ファイル
        body.append("--\(boundary)\r\n".data(using: .utf8)!)
        body.append("Content-Disposition: form-data; name=\"photo\"; filename=\"photo.jpg\"\r\n".data(using: .utf8)!)
        body.append("Content-Type: image/jpeg\r\n\r\n".data(using: .utf8)!)
        body.append(imageData)
        body.append("\r\n".data(using: .utf8)!)

        // 終了バウンダリ
        body.append("--\(boundary)--\r\n".data(using: .utf8)!)

        request.httpBody = body

        try await performRequest(request, successStatusCode: 201)
    }
}
