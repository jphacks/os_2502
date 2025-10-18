import Foundation

class APIService {
    static let shared = APIService()

    private let baseURL = URL(string: Configuration.shared.apiUrl)!

    private init() {}

    // MARK: - User APIs
    /// ユーザー作成
    func createUser(firebaseUID: String, name: String) async throws -> User {
        let url = URL(string: "\(baseURL)/api/users")!
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "firebase_uid": firebaseUID,
            "name": name
        ]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard httpResponse.statusCode == 201 else {
            throw APIError.httpError(statusCode: httpResponse.statusCode)
        }

        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        let user = try decoder.decode(User.self, from: data)
        return user
    }

    /// ユーザー削除
    func deleteUser(id: String) async throws {
        let url = URL(string: "\(baseURL)/api/users/\(id)")!
        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"

        let (_, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard httpResponse.statusCode == 204 else {
            throw APIError.httpError(statusCode: httpResponse.statusCode)
        }
    }
}

// MARK: - Models

extension APIService {
    struct User: Codable {
        let id: String
        let firebaseUid: String
        let name: String
        let username: String?
        let createdAt: String
        let updatedAt: String

        enum CodingKeys: String, CodingKey {
            case id
            case firebaseUid = "firebase_uid"
            case name
            case username
            case createdAt = "created_at"
            case updatedAt = "updated_at"
        }
    }
}

// MARK: - Errors

enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case httpError(statusCode: Int)
    case decodingError

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "無効なURLです"
        case .invalidResponse:
            return "無効なレスポンスです"
        case .httpError(let statusCode):
            return "HTTPエラー: \(statusCode)"
        case .decodingError:
            return "データの解析に失敗しました"
        }
    }
}
