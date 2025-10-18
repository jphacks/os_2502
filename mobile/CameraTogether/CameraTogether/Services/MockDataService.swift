import Foundation

/// モックデータを読み込むサービス
class MockDataService {
    static let shared = MockDataService()

    private init() {}

    /// JSONファイルからデータを読み込む
    func loadJSON<T: Decodable>(filename: String) throws -> T {
        // まずBundleから探す
        if let url = Bundle.main.url(forResource: filename, withExtension: "json", subdirectory: "Views/mock/db") {
            let data = try Data(contentsOf: url)
            let decoder = JSONDecoder()
            decoder.dateDecodingStrategy = .iso8601
            return try decoder.decode(T.self, from: data)
        }

        // Bundleに見つからない場合は埋め込みデータを使用
        let jsonString: String
        switch filename {
        case "groups":
            jsonString = Self.groupsJSON
        case "users":
            jsonString = Self.usersJSON
        default:
            throw MockDataError.fileNotFound(filename: filename)
        }

        guard let data = jsonString.data(using: .utf8) else {
            throw MockDataError.decodingFailed(error: NSError(domain: "MockDataService", code: -1))
        }

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601

        do {
            return try decoder.decode(T.self, from: data)
        } catch {
            throw MockDataError.decodingFailed(error: error)
        }
    }

    /// グループ一覧を取得
    func getGroups() throws -> GroupListResponse {
        return try loadJSON(filename: "groups")
    }

    /// ユーザー一覧を取得
    func getUsers() throws -> UserListResponse {
        return try loadJSON(filename: "users")
    }

    // MARK: - 埋め込みモックデータ

    private static let groupsJSON = """
    {
      "groups": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440001",
          "owner_user_id": "550e8400-e29b-41d4-a716-446655440010",
          "name": "家族グループ",
          "group_type": "permanent",
          "status": "recruiting",
          "max_member": 100,
          "current_member_count": 4,
          "invitation_token": "550e8400-e29b-41d4-a716-446655440101",
          "finalized_at": null,
          "countdown_started_at": null,
          "expires_at": null,
          "created_at": "2025-10-15T10:00:00Z",
          "updated_at": "2025-10-15T10:00:00Z"
        },
        {
          "id": "550e8400-e29b-41d4-a716-446655440002",
          "owner_user_id": "550e8400-e29b-41d4-a716-446655440010",
          "name": "友達グループ",
          "group_type": "global_temporary",
          "status": "ready_check",
          "max_member": 6,
          "current_member_count": 6,
          "invitation_token": "550e8400-e29b-41d4-a716-446655440102",
          "finalized_at": "2025-10-17T15:30:00Z",
          "countdown_started_at": null,
          "expires_at": "2025-10-20T23:59:59Z",
          "created_at": "2025-10-17T14:00:00Z",
          "updated_at": "2025-10-17T15:30:00Z"
        },
        {
          "id": "550e8400-e29b-41d4-a716-446655440003",
          "owner_user_id": "550e8400-e29b-41d4-a716-446655440011",
          "name": "旅行の思い出",
          "group_type": "global_temporary",
          "status": "completed",
          "max_member": 5,
          "current_member_count": 5,
          "invitation_token": "550e8400-e29b-41d4-a716-446655440103",
          "finalized_at": "2025-10-16T09:00:00Z",
          "countdown_started_at": "2025-10-16T09:05:00Z",
          "expires_at": "2025-10-19T23:59:59Z",
          "created_at": "2025-10-16T08:00:00Z",
          "updated_at": "2025-10-16T12:00:00Z"
        }
      ],
      "total_count": 3
    }
    """

    private static let usersJSON = """
    {
      "users": [
        {
          "id": "550e8400-e29b-41d4-a716-446655440010",
          "firebase_uid": "test_firebase_uid_001",
          "name": "山田太郎",
          "username": "yamada_taro",
          "created_at": "2025-10-15T10:00:00Z",
          "updated_at": "2025-10-15T10:00:00Z"
        },
        {
          "id": "550e8400-e29b-41d4-a716-446655440011",
          "firebase_uid": "test_firebase_uid_002",
          "name": "鈴木花子",
          "username": "suzuki_hanako",
          "created_at": "2025-10-15T11:00:00Z",
          "updated_at": "2025-10-15T11:00:00Z"
        }
      ],
      "limit": 10,
      "offset": 0,
      "count": 2
    }
    """
}

// MARK: - Errors

enum MockDataError: LocalizedError {
    case fileNotFound(filename: String)
    case decodingFailed(error: Error)

    var errorDescription: String? {
        switch self {
        case .fileNotFound(let filename):
            return "モックデータファイルが見つかりません: \(filename)"
        case .decodingFailed(let error):
            return "モックデータのデコードに失敗しました: \(error.localizedDescription)"
        }
    }
}
