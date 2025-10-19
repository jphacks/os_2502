import Foundation

/// テンプレート関連のAPI通信を管理するサービス
class TemplateAPIService: APIServiceBase {
    static let shared = TemplateAPIService()

    private override init() {
        super.init()
    }

    // MARK: - Template APIs

    /// すべてのテンプレートを取得
    /// - Returns: テンプレート一覧
    func getAllTemplates() async throws -> [CollageTemplate] {
        let url = baseURL.appendingPathComponent("template-data")
        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        let response = try await performRequest(request, expecting: TemplateDataResponse.self)
        return response.templates
    }

    /// 写真枚数でフィルタしてテンプレートを取得
    /// - Parameter photoCount: 写真枚数
    /// - Returns: フィルタされたテンプレート一覧
    func getTemplates(photoCount: Int) async throws -> [CollageTemplate] {
        var components = URLComponents(
            url: baseURL.appendingPathComponent("template-data/filter"),
            resolvingAgainstBaseURL: false
        )!
        components.queryItems = [URLQueryItem(name: "photo_count", value: String(photoCount))]

        guard let url = components.url else {
            throw APIError.invalidURL
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        let response = try await performRequest(request, expecting: TemplateDataResponse.self)
        return response.templates
    }

    /// IDでテンプレートを取得
    /// - Parameter id: テンプレートID
    /// - Returns: テンプレート
    func getTemplate(id: String) async throws -> CollageTemplate {
        let url = baseURL.appendingPathComponent("template-data").appendingPathComponent(id)
        var request = URLRequest(url: url)
        request.httpMethod = "GET"

        return try await performRequest(request, expecting: CollageTemplate.self)
    }
}
