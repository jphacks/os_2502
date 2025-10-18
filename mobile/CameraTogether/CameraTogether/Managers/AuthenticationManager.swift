import FirebaseAuth
import Foundation

@MainActor
@Observable
class AuthenticationManager {
    var user: FirebaseAuth.User?
    var backendUser: APIService.User?
    var isAuthenticated: Bool {
        user != nil
    }

    init() {
        checkAuth()
    }

    func checkAuth() {
        user = Auth.auth().currentUser
    }

    /// バックエンドにユーザーを作成
    private func createBackendUser(firebaseUID: String) async {
        do {
            let displayName = user?.displayName ?? "ユーザー"
            let newUser = try await APIService.shared.createUser(
                firebaseUID: firebaseUID,
                name: displayName
            )
            backendUser = newUser
        } catch {
            // 409エラー（既に存在）の場合は無視
            if let apiError = error as? APIError,
               case .httpError(let statusCode) = apiError,
               statusCode == 409 {
            }
        }
    }

    // MARK: - Apple Sign In
    func signInWithApple(idToken: String, nonce: String) async throws {
        let credential = OAuthProvider.credential(
            providerID: .apple,
            idToken: idToken,
            rawNonce: nonce
        )
        let result = try await Auth.auth().signIn(with: credential)
        user = result.user

        // バックエンドユーザーを作成（既に存在する場合は409エラーで無視）
        await createBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Google Sign In
    func signInWithGoogle(idToken: String, accessToken: String) async throws {
        let credential = GoogleAuthProvider.credential(
            withIDToken: idToken,
            accessToken: accessToken
        )
        let result = try await Auth.auth().signIn(with: credential)
        user = result.user

        // バックエンドユーザーを作成（既に存在する場合は409エラーで無視）
        await createBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Email/Password Sign In
    func signInWithEmail(email: String, password: String) async throws {
        let result = try await Auth.auth().signIn(withEmail: email, password: password)
        user = result.user

        // バックエンドユーザーを作成（既に存在する場合は409エラーで無視）
        await createBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Email/Password Sign Up
    func signUpWithEmail(email: String, password: String) async throws {
        let result = try await Auth.auth().createUser(withEmail: email, password: password)
        user = result.user

        // バックエンドユーザーを作成
        await createBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Sign Out
    func signOut() throws {
        try Auth.auth().signOut()
        user = nil
        backendUser = nil
    }

    // MARK: - Password Reset
    func resetPassword(email: String) async throws {
        try await Auth.auth().sendPasswordReset(withEmail: email)
    }

    // MARK: - Delete Account
    func deleteAccount() async throws {
        guard let currentUser = Auth.auth().currentUser else {
            throw AuthenticationError.noUser
        }

        // バックエンドのユーザーを削除
        if let backendUserId = backendUser?.id {
            try await APIService.shared.deleteUser(id: backendUserId)
        }

        // Firebaseのユーザーを削除
        try await currentUser.delete()

        user = nil
        backendUser = nil
    }
}

enum AuthenticationError: LocalizedError {
    case noUser

    var errorDescription: String? {
        switch self {
        case .noUser:
            return "ユーザーが見つかりません"
        }
    }
}
