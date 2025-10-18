import AuthenticationServices
import SwiftUI
import GoogleSignIn

struct LoginView: View {
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors
    var authManager: AuthenticationManager
    @State private var showEmailLogin = false
    @State private var currentNonce: String?
    @State private var errorMessage: String?
    @State private var isLoading = false

    var body: some View {
        ZStack {
            appColors.backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 0) {
                Spacer()

                puzzleIcon
                    .padding(.bottom, 140)

                VStack(spacing: 11) {
                    appleSignInButton
                    googleSignInButton
                    emailSignInButton
                }
                .padding(.horizontal, 24)
                .padding(.bottom, 84)

                if let errorMessage = errorMessage {
                    Text(errorMessage)
                        .font(.caption)
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)
                        .padding()
                }
            }
        }
        .sheet(isPresented: $showEmailLogin) {
            EmailLoginView(authManager: authManager)
        }
    }

    private var puzzleIcon: some View {
        Image("login-peace-logo")
            .resizable()
            .aspectRatio(contentMode: .fit)
            .frame(width: 217, height: 245)
    }

    private var appleSignInButton: some View {
        AppleSignInButton(
            onRequest: { request, nonce in
                currentNonce = nonce
            },
            onCompletion: { result in
                Task {
                    await handleAppleSignIn(result: result)
                }
            }
        )
        .frame(height: 42)
        .cornerRadius(45)
        .disabled(isLoading)
    }

    private var googleSignInButton: some View {
        Button {
            Task {
                await handleGoogleSignIn()
            }
        } label: {
            HStack(spacing: 8) {
                Image("material-icon-theme_google")
                    .frame(width: 16, height: 16)

                Text("Google で続ける")
                    .font(.system(size: 14.8, weight: .regular))
                    .foregroundColor(.black)
            }
            .frame(maxWidth: .infinity)
            .frame(height: 42)
            .background(Color.white)
            .cornerRadius(45)
        }
        .disabled(isLoading)
    }

    private var emailSignInButton: some View {
        Button {
            errorMessage = nil
            showEmailLogin = true
        } label: {
            HStack(spacing: 8) {
                Image(systemName: "envelope.fill")
                    .font(.system(size: 16, weight: .regular))
                    .foregroundColor(.black)

                Text("メールアドレスで続ける")
                    .font(.system(size: 14.8, weight: .regular))
                    .foregroundColor(.black)
            }
            .frame(maxWidth: .infinity)
            .frame(height: 42)
            .background(Color.white)
            .cornerRadius(45)
        }
        .disabled(isLoading)
    }

    private func handleAppleSignIn(result: Result<ASAuthorization, Error>) async {
        isLoading = true
        errorMessage = nil

        do {
            switch result {
            case .success(let authorization):
                guard
                    let appleIDCredential = authorization.credential
                        as? ASAuthorizationAppleIDCredential,
                    let nonce = currentNonce,
                    let appleIDToken = appleIDCredential.identityToken,
                    let idTokenString = String(data: appleIDToken, encoding: .utf8)
                else {
                    throw AuthError.invalidCredential
                }

                try await authManager.signInWithApple(idToken: idTokenString, nonce: nonce)

            case .failure(let error):
                throw error
            }
        } catch {
            // ユーザーがキャンセルした場合はエラーを表示しない
            if let authError = error as? ASAuthorizationError,
               authError.code == .canceled {
                // キャンセルは無視
            } else {
                errorMessage = error.localizedDescription
            }
        }

        isLoading = false
    }

    private func handleGoogleSignIn() async {
        isLoading = true
        errorMessage = nil

        do {
            let tokens = try await GoogleSignInHelper.signIn()
            try await authManager.signInWithGoogle(
                idToken: tokens.idToken, accessToken: tokens.accessToken)
        } catch {
            // ユーザーがキャンセルした場合はエラーを表示しない
            if let signInError = error as? GIDSignInError,
               signInError.code == .canceled {
                // キャンセルは無視
            } else {
                errorMessage = error.localizedDescription
            }
        }

        isLoading = false
    }
}

enum AuthError: LocalizedError {
    case invalidCredential

    var errorDescription: String? {
        switch self {
        case .invalidCredential:
            return "認証情報が無効です"
        }
    }
}

#Preview {
    LoginView(authManager: AuthenticationManager())
}
