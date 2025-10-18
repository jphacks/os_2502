import FirebaseAuth
import FirebaseCore
import GoogleSignIn
import SwiftUI

struct GoogleSignInHelper {
    static func signIn() async throws -> (idToken: String, accessToken: String) {
        guard let clientID = FirebaseApp.app()?.options.clientID else {
            throw GoogleSignInError.noClientID
        }

        let config = GIDConfiguration(clientID: clientID)
        GIDSignIn.sharedInstance.configuration = config

        guard
            let windowScene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
            let rootViewController = windowScene.windows.first?.rootViewController
        else {
            throw GoogleSignInError.noRootViewController
        }

        let result = try await GIDSignIn.sharedInstance.signIn(
            withPresenting: rootViewController)

        guard let idToken = result.user.idToken?.tokenString else {
            throw GoogleSignInError.noIDToken
        }

        let accessToken = result.user.accessToken.tokenString

        return (idToken: idToken, accessToken: accessToken)
    }
}

enum GoogleSignInError: LocalizedError {
    case noClientID
    case noRootViewController
    case noIDToken

    var errorDescription: String? {
        switch self {
        case .noClientID:
            return "Firebase Client IDが見つかりません"
        case .noRootViewController:
            return "ルートビューコントローラーが見つかりません"
        case .noIDToken:
            return "IDトークンが見つかりません"
        }
    }
}
