import Foundation

/// API通信の基底クラス
class APIServiceBase {
    /// ベースURL（/api付き）
    let baseURL = URL(string: Configuration.shared.apiUrl)!.appendingPathComponent("api").appendingPathComponent("api")

    /// JSONエンコーダー
    let encoder: JSONEncoder = {
        let encoder = JSONEncoder()
        encoder.dateEncodingStrategy = .iso8601
        return encoder
    }()

    /// JSONデコーダー
    let decoder: JSONDecoder = {
        let decoder = JSONDecoder()
        decoder.dateDecodingStrategy = .iso8601
        return decoder
    }()

    init() {}

    /// HTTPリクエストを実行
    func performRequest<T: Decodable>(
        _ request: URLRequest,
        expecting type: T.Type,
        successStatusCode: Int = 200
    ) async throws -> T {
        // リクエスト詳細をログ出力
        do {
            let (data, response) = try await URLSession.shared.data(for: request)

            guard let httpResponse = response as? HTTPURLResponse else {
                throw APIError.invalidResponse
            }

            guard httpResponse.statusCode == successStatusCode else {
                // エラー時のレスポンスボディをログ出力
                if let errorBody = String(data: data, encoding: .utf8) {
                    print("❌ API Error Response: \(errorBody)")
                }
                throw APIError.httpError(
                    statusCode: httpResponse.statusCode, message: String(data: data, encoding: .utf8))
            }

            do {
                return try decoder.decode(T.self, from: data)
            } catch {
                throw APIError.decodingError
            }
        } catch let error as APIError {
            throw error
        } catch {
            throw error
        }
    }

    /// HTTPリクエストを実行（レスポンスなし）
    func performRequest(
        _ request: URLRequest,
        successStatusCode: Int = 200
    ) async throws {
        let (data, response) = try await URLSession.shared.data(for: request)

        guard let httpResponse = response as? HTTPURLResponse else {
            throw APIError.invalidResponse
        }

        guard httpResponse.statusCode == successStatusCode else {
            // エラー時のレスポンスボディをログ出力
            if let errorBody = String(data: data, encoding: .utf8) {
                print(
                    "❌ API Error [\(httpResponse.statusCode)]: \(request.url?.absoluteString ?? "")"
                )
                print("Response: \(errorBody)")
            }
            throw APIError.httpError(
                statusCode: httpResponse.statusCode, message: String(data: data, encoding: .utf8))
        }
    }
}

// MARK: - Errors

enum APIError: LocalizedError {
    case invalidURL
    case invalidResponse
    case httpError(statusCode: Int, message: String? = nil)
    case decodingError

    var errorDescription: String? {
        switch self {
        case .invalidURL:
            return "無効なURLです"
        case .invalidResponse:
            return "無効なレスポンスです"
        case .httpError(let statusCode, let message):
            if let message = message {
                return "HTTPエラー\(statusCode): \(message)"
            }
            return "HTTPエラー\(statusCode)"
        case .decodingError:
            return "データの解析に失敗しました"
        }
    }
}
