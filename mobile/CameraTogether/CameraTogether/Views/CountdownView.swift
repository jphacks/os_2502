import AVFoundation
import SwiftUI

struct CountdownView: View {
    @Bindable var viewModel: CollageGroupViewModel
    @StateObject private var cameraService = CameraService()
    @State private var countdown = 10
    @State private var isCountingDown = true
    @State private var capturedImage: UIImage?
    @State private var showingTemplateSelection = false
    @State private var showingError = false
    @State private var errorMessage = ""

    var body: some View {
        ZStack {
            CameraPreviewView(cameraService: cameraService)
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

                    ProgressView(value: Double(10 - countdown), total: 10)
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
                            showingTemplateSelection = true
                        } label: {
                            Text("テンプレート選択へ")
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
        .alert("エラー", isPresented: $showingError) {
            Button("OK", role: .cancel) {}
        } message: {
            Text(errorMessage)
        }
        .task {
            await setupCamera()
        }
        .onAppear {
            startCountdown()
        }
        .onDisappear {
            cameraService.stopSession()
        }
        .navigationDestination(isPresented: $showingTemplateSelection) {
            TemplateSelectionView(viewModel: viewModel, capturedImage: capturedImage ?? UIImage())
        }
    }

    private func setupCamera() async {
        do {
            let authorized = await cameraService.checkAuthorization()
            if !authorized {
                await MainActor.run {
                    errorMessage = "カメラへのアクセスが許可されていません"
                    showingError = true
                }
                return
            }

            try cameraService.setupSession()
            cameraService.startSession()
        } catch {
            await MainActor.run {
                errorMessage = "カメラの起動に失敗しました: \(error.localizedDescription)"
                showingError = true
            }
        }
    }

    private func startCountdown() {
        Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { timer in
            if countdown > 0 {
                countdown -= 1
            } else {
                timer.invalidate()
                Task {
                    await capturePhoto()
                }
            }
        }
    }

    private func capturePhoto() async {
        do {
            let image = try await cameraService.capturePhoto()
            await MainActor.run {
                isCountingDown = false
                capturedImage = image
                viewModel.completeSession()
            }
        } catch {
            await MainActor.run {
                errorMessage = "撮影に失敗しました: \(error.localizedDescription)"
                showingError = true
                isCountingDown = false
            }
        }
    }
}

struct CameraPreviewView: UIViewRepresentable {
    @ObservedObject var cameraService: CameraService

    func makeUIView(context: Context) -> UIView {
        let view = UIView(frame: .zero)
        view.backgroundColor = .black

        let previewLayer = cameraService.getPreviewLayer()
        previewLayer.frame = view.bounds
        view.layer.addSublayer(previewLayer)

        // レイヤーのフレームを更新するためのコンテキストを保存
        context.coordinator.previewLayer = previewLayer

        return view
    }

    func updateUIView(_ uiView: UIView, context: Context) {
        // プレビューレイヤーのフレームを更新
        if let previewLayer = context.coordinator.previewLayer {
            DispatchQueue.main.async {
                previewLayer.frame = uiView.bounds
            }
        }
    }

    func makeCoordinator() -> Coordinator {
        Coordinator()
    }

    class Coordinator {
        var previewLayer: AVCaptureVideoPreviewLayer?
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
