import Foundation

/// アプリ全体で使用する開発フラグと設定
struct AppConfig {
    /// モックデータを使用するかどうか
    static var useMockData: Bool {
        #if DEBUG
        return true  // デバッグビルドではデフォルトでモックデータを使用
        #else
        return false
        #endif
    }

    /// APIのベースURL
    static var apiBaseURL: String {
        return Configuration.shared.apiUrl
    }

    /// ログ出力を有効にするか
    static var enableLogging: Bool {
        #if DEBUG
        return true
        #else
        return false
        #endif
    }
}
