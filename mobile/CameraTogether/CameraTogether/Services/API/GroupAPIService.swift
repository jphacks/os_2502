import Foundation

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
    func getGroups(ownerUserId: String? = nil, limit: Int = 10, offset: Int = 0) async throws -> [APIGroup] {
        // モックデータを使用する場合
        if AppConfig.useMockData {
            let response = try MockDataService.shared.getGroups()
            return response.groups
        }

        // 実際のAPI呼び出し
        var components = URLComponents(url: baseURL.appendingPathComponent("groups"), resolvingAgainstBaseURL: false)!
        var queryItems: [URLQueryItem] = [
            URLQueryItem(name: "limit", value: "\(limit)"),
            URLQueryItem(name: "offset", value: "\(offset)")
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
    func createGroup(ownerUserId: String, name: String, groupType: String = "global_temporary") async throws -> APIGroup {
        let url = baseURL.appendingPathComponent("groups")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "owner_user_id": ownerUserId,
            "name": name,
            "group_type": groupType
        ]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        return try await performRequest(request, expecting: APIGroup.self, successStatusCode: 201)
    }

    /// グループ削除
    /// - Parameters:
    ///   - id: グループID
    ///   - userId: ユーザーID（オーナー確認用）
    func deleteGroup(id: String, userId: String) async throws {
        var components = URLComponents(url: baseURL.appendingPathComponent("groups").appendingPathComponent(id), resolvingAgainstBaseURL: false)!
        components.queryItems = [URLQueryItem(name: "user_id", value: userId)]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"

        try await performRequest(request, successStatusCode: 200)
    }
}
