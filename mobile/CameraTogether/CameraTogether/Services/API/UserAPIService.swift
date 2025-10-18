import Foundation

/// ユーザー関連のAPI通信を管理するサービス
class UserAPIService: APIServiceBase {
    static let shared = UserAPIService()

    private override init() {
        super.init()
    }

    // MARK: - User APIs

    /// ユーザー作成
    /// - Parameters:
    ///   - firebaseUID: Firebase Authentication UID
    ///   - name: ユーザー名
    /// - Returns: 作成されたユーザー
    func createUser(firebaseUID: String, name: String) async throws -> User {
        let url = baseURL.appendingPathComponent("users")
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

        let body: [String: Any] = [
            "firebase_uid": firebaseUID,
            "name": name,
        ]
        request.httpBody = try JSONSerialization.data(withJSONObject: body)

        return try await performRequest(request, expecting: User.self, successStatusCode: 201)
    }

    /// Firebase UIDでユーザー取得
    /// - Parameter firebaseUID: Firebase Authentication UID
    /// - Returns: ユーザー情報
    func getUserByFirebaseUID(firebaseUID: String) async throws -> User {
        var components = URLComponents(
            url: baseURL.appendingPathComponent("users/firebase"), resolvingAgainstBaseURL: false)!
        components.queryItems = [URLQueryItem(name: "firebase_uid", value: firebaseUID)]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        return try await performRequest(request, expecting: User.self)
    }

    /// ユーザー削除
    /// - Parameter id: ユーザーID
    func deleteUser(id: String) async throws {
        let url = baseURL.appendingPathComponent("users").appendingPathComponent(id)
        var request = URLRequest(url: url)
        request.httpMethod = "DELETE"

        try await performRequest(request, successStatusCode: 204)
    }
}
