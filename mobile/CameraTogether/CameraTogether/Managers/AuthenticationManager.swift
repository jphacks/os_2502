import FirebaseAuth
import Foundation

@MainActor
@Observable
class AuthenticationManager {
    var user: User?
    var isAuthenticated: Bool {
        user != nil
    }

    init() {
        checkAuth()
    }

    func checkAuth() {
        user = Auth.auth().currentUser
    }

    // Apple Sign In
    func signInWithApple(idToken: String, nonce: String) async throws {
        let credential = OAuthProvider.credential(
            providerID: .apple,
            idToken: idToken,
            rawNonce: nonce
        )
        let result = try await Auth.auth().signIn(with: credential)
        user = result.user
    }

    // Google Sign In
    func signInWithGoogle(idToken: String, accessToken: String) async throws {
        let credential = GoogleAuthProvider.credential(
            withIDToken: idToken,
            accessToken: accessToken
        )
        let result = try await Auth.auth().signIn(with: credential)
        user = result.user
    }

    // Email/Password Sign In
    func signInWithEmail(email: String, password: String) async throws {
        let result = try await Auth.auth().signIn(withEmail: email, password: password)
        user = result.user
    }

    // Email/Password Sign Up
    func signUpWithEmail(email: String, password: String) async throws {
        let result = try await Auth.auth().createUser(withEmail: email, password: password)
        user = result.user
    }

    // Sign Out
    func signOut() throws {
        try Auth.auth().signOut()
        user = nil
    }

    // Password Reset
    func resetPassword(email: String) async throws {
        try await Auth.auth().sendPasswordReset(withEmail: email)
    }

    // Delete Account
    func deleteAccount() async throws {
        guard let currentUser = Auth.auth().currentUser else {
            throw AuthenticationError.noUser
        }
        try await currentUser.delete()
        user = nil
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
