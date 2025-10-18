import AVFoundation
import SwiftUI

struct CountdownView: View {
    @Bindable var viewModel: CollageGroupViewModel
    @State private var countdown = 20
    @State private var isCountingDown = true
    @State private var capturedImage: UIImage?
    @State private var showingResult = false

    var body: some View {
        ZStack {
            CameraPreviewView()
                .ignoresSafeArea()

            VStack {
                if isCountingDown {
                    Spacer()

                    Text("\(countdown)")
                        .font(.system(size: 120, weight: .bold))
                        .foregroundColor(.white)
                        .shadow(color: .black.opacity(0.5), radius: 10)

                    Spacer()

                    Text("撮影まで")
                        .font(.title)
                        .foregroundColor(.white)
                        .shadow(color: .black.opacity(0.5), radius: 5)

                    ProgressView(value: Double(20 - countdown), total: 20)
                        .progressViewStyle(LinearProgressViewStyle(tint: .white))
                        .scaleEffect(x: 1, y: 4, anchor: .center)
                        .padding(.horizontal, 40)
                        .padding(.bottom, 60)
                } else if capturedImage != nil {
                    VStack {
                        Spacer()

                        Text("撮影完了")
                            .font(.largeTitle)
                            .fontWeight(.bold)
                            .foregroundColor(.white)
                            .shadow(color: .black.opacity(0.5), radius: 10)

                        Spacer()

                        Button {
                            showingResult = true
                        } label: {
                            Text("結果を見る")
                                .font(.title2)
                                .fontWeight(.semibold)
                                .foregroundColor(.white)
                                .frame(maxWidth: .infinity)
                                .padding()
                                .background(Color.blue)
                                .cornerRadius(12)
                        }
                        .padding(.horizontal, 40)
                        .padding(.bottom, 60)
                    }
                }
            }
        }
        .navigationBarBackButtonHidden(true)
        .onAppear {
            startCountdown()
        }
        .navigationDestination(isPresented: $showingResult) {
            ResultView(viewModel: viewModel, capturedImage: capturedImage)
        }
    }

    private func startCountdown() {
        Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { timer in
            if countdown > 0 {
                countdown -= 1
            } else {
                timer.invalidate()
                capturePhoto()
            }
        }
    }

    private func capturePhoto() {
        // 実際のカメラ撮影の代わりにダミー画像を使用
        // 実装時はCameraViewから撮影機能を統合
        isCountingDown = false

        // ダミーの撮影処理
        DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
            capturedImage = UIImage(systemName: "photo")
            viewModel.completeSession()
        }
    }
}

struct CameraPreviewView: UIViewRepresentable {
    func makeUIView(context: Context) -> UIView {
        let view = UIView(frame: .zero)
        view.backgroundColor = .black

        // カメラプレビュー用のレイヤー
        // 実装時はAVCaptureSessionを使用
        let previewView = UIView()
        previewView.backgroundColor = .darkGray
        previewView.translatesAutoresizingMaskIntoConstraints = false
        view.addSubview(previewView)

        NSLayoutConstraint.activate([
            previewView.topAnchor.constraint(equalTo: view.topAnchor),
            previewView.bottomAnchor.constraint(equalTo: view.bottomAnchor),
            previewView.leadingAnchor.constraint(equalTo: view.leadingAnchor),
            previewView.trailingAnchor.constraint(equalTo: view.trailingAnchor),
        ])

        return view
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        // カメラプレビューの更新処理
    }
}

#Preview {
    let authManager = AuthenticationManager()
    let viewModel = CollageGroupViewModel(authManager: authManager)
    let _ = viewModel.createGroupLocal(type: .temporaryLocal, maxMembers: 5)

    return NavigationStack {
        CountdownView(viewModel: viewModel)
    }
}
