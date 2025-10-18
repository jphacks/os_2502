import FirebaseAuth
import SwiftUI

struct EmailLoginView: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @State private var email = ""
    @State private var password = ""
    @State private var isSignUp = false
    @State private var errorMessage: String?
    @State private var isLoading = false
    let authManager: AuthenticationManager

    var body: some View {
        NavigationStack {
            ZStack {
                Color.black
                    .ignoresSafeArea()

                ScrollView {
                    VStack(spacing: 24) {
                        Spacer()
                            .frame(height: 60)

                        Text(isSignUp ? "アカウント作成" : "ログイン")
                            .font(.largeTitle)
                            .fontWeight(.bold)
                            .foregroundColor(.white)

                        VStack(spacing: 16) {
                            emailField
                            passwordField

                            if let errorMessage = errorMessage {
                                Text(errorMessage)
                                    .font(.caption)
                                    .foregroundColor(.red)
                                    .multilineTextAlignment(.center)
                            }
                        }
                        .padding(.horizontal, 24)

                        VStack(spacing: 12) {
                            submitButton

                            Button {
                                isSignUp.toggle()
                            } label: {
                                Text(isSignUp ? "すでにアカウントをお持ちの方" : "アカウントを作成")
                                    .font(.system(size: 14.8, weight: .regular))
                                    .foregroundColor(.white.opacity(0.7))
                            }
                        }
                        .padding(.horizontal, 24)

                        Spacer()
                    }
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button {
                        errorMessage = nil
                        dismiss()
                    } label: {
                        Image(systemName: "xmark.circle.fill")
                            .font(.title3)
                            .foregroundColor(.white)
                    }
                }
            }
            .toolbarBackground(.visible, for: .navigationBar)
            .toolbarBackground(Color.clear, for: .navigationBar)
        }
    }

    private var emailField: some View {
        TextField("", text: $email, prompt: Text("メールアドレス").foregroundColor(.white.opacity(0.5)))
            .textContentType(.emailAddress)
            .keyboardType(.emailAddress)
            .autocapitalization(.none)
            .foregroundColor(.white)
            .padding()
            .background(
                RoundedRectangle(cornerRadius: 12)
                    .fill(Color.white.opacity(0.1))
                    .overlay(
                        RoundedRectangle(cornerRadius: 12)
                            .stroke(Color.white.opacity(0.2), lineWidth: 1)
                    )
            )
    }

    private var passwordField: some View {
        SecureField(
            "", text: $password, prompt: Text("パスワード").foregroundColor(.white.opacity(0.5))
        )
        .textContentType(isSignUp ? .newPassword : .password)
        .foregroundColor(.white)
        .padding()
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(Color.white.opacity(0.1))
                .overlay(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(Color.white.opacity(0.2), lineWidth: 1)
                )
        )
    }

    private var submitButton: some View {
        Button {
            Task {
                await handleSubmit()
            }
        } label: {
            HStack {
                if isLoading {
                    ProgressView()
                        .progressViewStyle(CircularProgressViewStyle(tint: .black))
                } else {
                    Text(isSignUp ? "登録" : "ログイン")
                        .font(.system(size: 16, weight: .semibold))
                }
            }
            .frame(maxWidth: .infinity)
            .frame(height: 50)
            .background(Color.white)
            .foregroundColor(.black)
            .cornerRadius(25)
        }
        .disabled(isLoading || email.isEmpty || password.isEmpty)
        .opacity(isLoading || email.isEmpty || password.isEmpty ? 0.6 : 1.0)
    }

    private func handleSubmit() async {
        isLoading = true
        errorMessage = nil

        do {
            if isSignUp {
                try await authManager.signUpWithEmail(email: email, password: password)
            } else {
                try await authManager.signInWithEmail(email: email, password: password)
            }
            dismiss()
        } catch {
            errorMessage = friendlyErrorMessage(from: error)
        }

        isLoading = false
    }

    private func friendlyErrorMessage(from error: Error) -> String {
        if let authError = error as? AuthErrorCode {
            switch authError.code {
            case .invalidEmail:
                return "メールアドレスの形式が正しくありません"
            case .emailAlreadyInUse:
                return "このメールアドレスは既に使用されています"
            case .weakPassword:
                return "パスワードは6文字以上で設定してください"
            case .wrongPassword:
                return "メールアドレスまたはパスワードが間違っています"
            case .userNotFound:
                return "このメールアドレスは登録されていません"
            case .networkError:
                return "ネットワークエラーが発生しました"
            case .tooManyRequests:
                return "しばらく時間をおいてから再度お試しください"
            default:
                return error.localizedDescription
            }
        }
        return error.localizedDescription
    }
}

#Preview {
    @Previewable @State var authManager = AuthenticationManager()
    EmailLoginView(authManager: authManager)
}
