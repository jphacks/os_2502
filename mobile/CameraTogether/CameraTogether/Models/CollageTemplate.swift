import Foundation

/// コラージュテンプレートのフレーム情報
struct CollageTemplateFrame: Codable, Identifiable {
    let id: Int
    let path: String  // SVG path データ

    enum CodingKeys: String, CodingKey {
        case id
        case path
    }
}

/// コラージュテンプレート
struct CollageTemplate: Codable, Identifiable {
    var id: String { name }  // nameをIDとして使用
    let name: String
    let photoCount: Int
    let viewBox: String
    let frames: [CollageTemplateFrame]

    enum CodingKeys: String, CodingKey {
        case name
        case photoCount = "photo_count"
        case viewBox
        case frames
    }
}

/// テンプレートAPIレスポンス
struct TemplateDataResponse: Codable {
    let templates: [CollageTemplate]
    let count: Int
    let photoCount: Int?

    enum CodingKeys: String, CodingKey {
        case templates
        case count
        case photoCount = "photo_count"
    }
}
