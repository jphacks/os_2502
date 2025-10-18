import AuthenticationServices
import CryptoKit
import SwiftUI

struct AppleSignInHelper {
    static func randomNonceString(length: Int = 32) -> String {
        precondition(length > 0)
        var randomBytes = [UInt8](repeating: 0, count: length)
        let errorCode = SecRandomCopyBytes(kSecRandomDefault, randomBytes.count, &randomBytes)
        if errorCode != errSecSuccess {
            fatalError(
                "Unable to generate nonce. SecRandomCopyBytes failed with OSStatus \(errorCode)")
        }

        let charset: [Character] =
            Array("0123456789ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvwxyz-_")

        let nonce = randomBytes.map { byte in
            charset[Int(byte) % charset.count]
        }

        return String(nonce)
    }

    static func sha256(_ input: String) -> String {
        let inputData = Data(input.utf8)
        let hashedData = SHA256.hash(data: inputData)
        let hashString = hashedData.compactMap {
            String(format: "%02x", $0)
        }.joined()

        return hashString
    }
}

struct AppleSignInButton: UIViewRepresentable {
    @Environment(\.colorScheme) var colorScheme
    var onRequest: (ASAuthorizationAppleIDRequest, String) -> Void
    var onCompletion: (Result<ASAuthorization, Error>) -> Void

    func makeUIView(context: Context) -> ASAuthorizationAppleIDButton {
        let button = ASAuthorizationAppleIDButton(
            authorizationButtonType: .continue,
            authorizationButtonStyle: colorScheme == .dark ? .white : .black
        )
        button.addTarget(
            context.coordinator,
            action: #selector(Coordinator.didTapButton),
            for: .touchUpInside
        )
        return button
    }

    func updateUIView(_ uiView: ASAuthorizationAppleIDButton, context: Context) {
    }

    func makeCoordinator() -> Coordinator {
        Coordinator(self)
    }

    class Coordinator: NSObject, ASAuthorizationControllerDelegate,
        ASAuthorizationControllerPresentationContextProviding
    {
        let parent: AppleSignInButton
        var currentNonce: String?

        init(_ parent: AppleSignInButton) {
            self.parent = parent
        }

        @objc func didTapButton() {
            let nonce = AppleSignInHelper.randomNonceString()
            currentNonce = nonce
            let appleIDProvider = ASAuthorizationAppleIDProvider()
            let request = appleIDProvider.createRequest()
            request.requestedScopes = [.fullName, .email]
            request.nonce = AppleSignInHelper.sha256(nonce)

            parent.onRequest(request, nonce)

            let authorizationController = ASAuthorizationController(authorizationRequests: [
                request
            ])
            authorizationController.delegate = self
            authorizationController.presentationContextProvider = self
            authorizationController.performRequests()
        }

        func authorizationController(
            controller: ASAuthorizationController,
            didCompleteWithAuthorization authorization: ASAuthorization
        ) {
            parent.onCompletion(.success(authorization))
        }

        func authorizationController(
            controller: ASAuthorizationController, didCompleteWithError error: Error
        ) {
            parent.onCompletion(.failure(error))
        }

        func presentationAnchor(for controller: ASAuthorizationController) -> ASPresentationAnchor {
            guard
                let windowScene = UIApplication.shared.connectedScenes
                    .compactMap({ $0 as? UIWindowScene })
                    .first
            else {
                fatalError("Unable to get window scene for Apple Sign In")
            }

            // iOS 26.0以降は新しいイニシャライザを使用
            if #available(iOS 26.0, *) {
                return ASPresentationAnchor(windowScene: windowScene)
            } else {
                // iOS 25.0以前は従来の方法
                guard let window = windowScene.windows.first else {
                    fatalError("Unable to get window for Apple Sign In")
                }
                return window
            }
        }
    }
}
