import FirebaseAuth
import Foundation

@MainActor
@Observable
class AuthenticationManager {
    var user: FirebaseAuth.User?
    var backendUser: User?
    var isAuthenticated: Bool {
        user != nil
    }

    init() {
        checkAuth()
    }

    func checkAuth() {
        user = Auth.auth().currentUser

        // 既にFirebase認証されている場合、バックエンドユーザー情報を取得
        if let firebaseUser = user {
            Task {
                await createOrFetchBackendUser(firebaseUID: firebaseUser.uid)
            }
        }
    }

    /// バックエンドにユーザーを作成または取得
    private func createOrFetchBackendUser(firebaseUID: String) async {
        do {
            let displayName = user?.displayName ?? "ユーザー"
            let newUser = try await UserAPIService.shared.createUser(
                firebaseUID: firebaseUID,
                name: displayName
            )
            backendUser = newUser
        } catch {
            // 409エラー（既に存在）の場合は既存ユーザーを取得
            if let apiError = error as? APIError,
                case .httpError(let statusCode) = apiError,
                statusCode == 409
            {
                do {
                    let existingUser = try await UserAPIService.shared.getUserByFirebaseUID(
                        firebaseUID: firebaseUID)
                    backendUser = existingUser
                } catch {
                    print("既存ユーザーの取得に失敗: \(error)")
                }
            } else {
                print("バックエンドユーザー作成に失敗: \(error)")
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

        // バックエンドユーザーを作成または取得
        await createOrFetchBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Google Sign In
    func signInWithGoogle(idToken: String, accessToken: String) async throws {
        let credential = GoogleAuthProvider.credential(
            withIDToken: idToken,
            accessToken: accessToken
        )
        let result = try await Auth.auth().signIn(with: credential)
        user = result.user

        // バックエンドユーザーを作成または取得
        await createOrFetchBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Email/Password Sign In
    func signInWithEmail(email: String, password: String) async throws {
        let result = try await Auth.auth().signIn(withEmail: email, password: password)
        user = result.user

        // バックエンドユーザーを作成または取得
        await createOrFetchBackendUser(firebaseUID: result.user.uid)
    }

    // MARK: - Email/Password Sign Up
    func signUpWithEmail(email: String, password: String) async throws {
        let result = try await Auth.auth().createUser(withEmail: email, password: password)
        user = result.user

        // バックエンドユーザーを作成
        await createOrFetchBackendUser(firebaseUID: result.user.uid)
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
            try await UserAPIService.shared.deleteUser(id: backendUserId)
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
