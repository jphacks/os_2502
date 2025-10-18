import Foundation

/// モックデータを読み込むサービス
class MockDataService {
    static let shared = MockDataService()

    private init() {}

    /// JSONファイルからデータを読み込む
    func loadJSON<T: Decodable>(filename: String) throws -> T {
        guard let url = Bundle.main.url(forResource: filename, withExtension: "json", subdirectory: "Views/mock/db") else {
            throw MockDataError.fileNotFound(filename: filename)
        }

        let data = try Data(contentsOf: url)
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
