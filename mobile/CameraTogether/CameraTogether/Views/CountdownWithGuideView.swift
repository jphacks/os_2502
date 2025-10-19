import AVFoundation
import SwiftUI

/// テンプレートのフレームガイド付きカウントダウン画面
struct CountdownWithGuideView: View {
    @Bindable var viewModel: CollageGroupViewModel
    let template: CollageTemplate
    let myFrameIndex: Int

    @State private var cameraService: CameraService?
    @State private var countdown = 10
    @State private var isCountingDown = true
    @State private var capturedImage: UIImage?
    @State private var showingResult = false
    @State private var showingError = false
    @State private var errorMessage = ""
    @State private var countdownTimer: Timer?

    var body: some View {
        ZStack {
            if let service = cameraService {
                CameraPreviewWrapper(cameraService: service)
                    .ignoresSafeArea()
            } else {
                Color.black.ignoresSafeArea()
            }

            // フレームガイドオーバーレイ
            if isCountingDown, myFrameIndex < template.frames.count {
                FrameGuideOverlay(
                    frame: template.frames[myFrameIndex],
                    viewBox: template.viewBox
                )
                .ignoresSafeArea()
            }

            VStack {
                if isCountingDown {
                    Spacer()

                    VStack(spacing: 8) {
                        Text("あなたの担当パート")
                            .font(.caption)
                            .foregroundColor(.white.opacity(0.8))
                            .shadow(color: .black.opacity(0.5), radius: 5)

                        Text("\(countdown)")
                            .font(.system(size: 120, weight: .bold))
                            .foregroundColor(.white)
                            .shadow(color: .black.opacity(0.5), radius: 10)
                    }

                    Spacer()

                    Text("白い枠に合わせて撮影してください")
                        .font(.subheadline)
                        .foregroundColor(.white)
                        .shadow(color: .black.opacity(0.5), radius: 5)
                        .padding(.bottom, 8)

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
                            showingResult = true
                        } label: {
                            Text("確認")
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
            print("CountdownWithGuideView appeared")
            startCountdown()
        }
        .onDisappear {
            print("CountdownWithGuideView disappeared")
            countdownTimer?.invalidate()
            cameraService?.stopSession()
        }
        .navigationDestination(isPresented: $showingResult) {
            // TODO: 撮影結果確認画面
            Text("撮影した写真: Frame \(myFrameIndex + 1)")
        }
    }

    private func setupCamera() async {
        print("CountdownWithGuideView: setupCamera() started")

        // CameraServiceを遅延初期化
        let service = CameraService()
        await MainActor.run {
            cameraService = service
        }
        print("CameraService instance created")

        do {
            print("Checking camera authorization...")
            let authorized = await service.checkAuthorization()
            print("Camera authorized: \(authorized)")

            if !authorized {
                print("Camera not authorized")
                await MainActor.run {
                    errorMessage = "カメラへのアクセスが許可されていません"
                    showingError = true
                }
                return
            }

            print("Setting up camera session...")
            try service.setupSession()
            print("Starting camera session...")
            service.startSession()
            print("Camera setup complete")
        } catch {
            print("setupCamera error: \(error)")
            await MainActor.run {
                errorMessage = "カメラの起動に失敗しました: \(error.localizedDescription)"
                showingError = true
            }
        }
    }

    private func startCountdown() {
        // サーバーから指定された撮影予定時刻を使用
        guard let scheduledTime = viewModel.currentGroup?.scheduledCaptureTime else {
            print("No scheduled capture time, using local countdown")
            startLocalCountdown()
            return
        }

        print("🕐 Server scheduled capture time: \(scheduledTime)")
        print("🕐 Current time: \(Date())")

        // サーバー時刻までの残り時間を計算
        let timeInterval = scheduledTime.timeIntervalSinceNow

        if timeInterval <= 0 {
            print("Scheduled time has already passed, capturing immediately")
            Task {
                await capturePhoto()
            }
            return
        }

        // 初期カウントダウン値を設定（切り上げ）
        countdown = Int(ceil(timeInterval))
        print("Starting countdown from \(countdown) seconds")

        // 高精度タイマーで毎秒更新
        countdownTimer = Timer.scheduledTimer(withTimeInterval: 0.1, repeats: true) { [self] timer in
            let remainingTime = scheduledTime.timeIntervalSinceNow

            if remainingTime <= 0 {
                // 撮影時刻になった
                timer.invalidate()
                countdown = 0
                Task {
                    await capturePhoto()
                }
            } else {
                // カウントダウン表示を更新
                let newCountdown = Int(ceil(remainingTime))
                if newCountdown != countdown {
                    countdown = newCountdown
                    print("Countdown: \(countdown)")
                }
            }
        }
    }

    /// ローカルカウントダウン（フォールバック用）
    private func startLocalCountdown() {
        countdown = 10
        countdownTimer = Timer.scheduledTimer(withTimeInterval: 1.0, repeats: true) { timer in
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
        guard let cameraService = cameraService else {
            await MainActor.run {
                errorMessage = "カメラが初期化されていません"
                showingError = true
                isCountingDown = false
            }
            return
        }

        do {
            print("📸 Capturing photo...")
            let image = try await cameraService.capturePhoto()
            print("Photo captured successfully")

            await MainActor.run {
                isCountingDown = false
                capturedImage = image
            }

            // 画像をサーバーにアップロード
            await uploadPhotoToServer(image: image)

            await MainActor.run {
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

    private func uploadPhotoToServer(image: UIImage) async {
        guard let groupId = viewModel.currentGroup?.id else {
            print("No group ID available for upload")
            return
        }

        let userId = viewModel.currentUserId
        print("Starting photo upload...")

        do {
            try await GroupAPIService.shared.uploadPhoto(
                groupId: groupId,
                userId: userId,
                image: image,
                frameIndex: myFrameIndex
            )
            print("Photo uploaded to server successfully")
        } catch {
            print("Failed to upload photo: \(error.localizedDescription)")
            await MainActor.run {
                errorMessage = "写真のアップロードに失敗しました: \(error.localizedDescription)"
                showingError = true
            }
        }
    }
}

/// フレームガイドオーバーレイ
struct FrameGuideOverlay: View {
    let frame: CollageTemplateFrame
    let viewBox: String

    var body: some View {
        GeometryReader { geometry in
            ZStack {
                // 暗いオーバーレイ（ガイド外）
                Color.black.opacity(0.5)

                // ガイド部分を切り抜き
                FrameGuidePath(pathString: frame.path, viewBox: viewBox)
                    .stroke(Color.white, lineWidth: 3)
                    .shadow(color: .white.opacity(0.5), radius: 10)

                // 内側を透明に
                FrameGuidePath(pathString: frame.path, viewBox: viewBox)
                    .blendMode(.destinationOut)
            }
            .compositingGroup()
        }
    }
}

/// フレームガイドパス
struct FrameGuidePath: Shape {
    let pathString: String
    let viewBox: String

    func path(in rect: CGRect) -> Path {
        let components = viewBox.split(separator: " ").compactMap { Double($0) }
        guard components.count == 4 else {
            return Path()
        }

        let viewBoxWidth = components[2]
        let viewBoxHeight = components[3]
        let scaleX = rect.width / viewBoxWidth
        let scaleY = rect.height / viewBoxHeight

        var path = Path()
        var currentPoint = CGPoint.zero
        let commands = pathString.uppercased()
        var index = commands.startIndex

        while index < commands.endIndex {
            let command = commands[index]
            index = commands.index(after: index)

            switch command {
            case "M":
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 2) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: coords[1] * scaleY)
                    path.move(to: currentPoint)
                }
            case "L":
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 2) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: coords[1] * scaleY)
                    path.addLine(to: currentPoint)
                }
            case "H":
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 1) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: currentPoint.y)
                    path.addLine(to: currentPoint)
                }
            case "V":
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 1) {
                    currentPoint = CGPoint(x: currentPoint.x, y: coords[0] * scaleY)
                    path.addLine(to: currentPoint)
                }
            case "Z":
                path.closeSubpath()
            default:
                break
            }
        }

        return path
    }

    private func parseCoordinates(from string: String, startingAt index: inout String.Index, count: Int) -> [CGFloat]? {
        var coords: [CGFloat] = []
        var numberString = ""

        while index < string.endIndex && coords.count < count {
            let char = string[index]

            if char.isNumber || char == "." || char == "-" {
                numberString.append(char)
                index = string.index(after: index)
            } else if !numberString.isEmpty {
                if let number = Double(numberString) {
                    coords.append(CGFloat(number))
                }
                numberString = ""

                if char.isWhitespace || char == "," {
                    index = string.index(after: index)
                }
            } else {
                if char.isWhitespace || char == "," {
                    index = string.index(after: index)
                } else {
                    break
                }
            }
        }

        if !numberString.isEmpty, let number = Double(numberString) {
            coords.append(CGFloat(number))
        }

        return coords.count == count ? coords : nil
    }
}

/// CameraPreviewのラッパー（オプショナル対応）
struct CameraPreviewWrapper: UIViewRepresentable {
    let cameraService: CameraService

    func makeUIView(context: Context) -> UIView {
        let view = UIView(frame: .zero)
        view.backgroundColor = .black

        let previewLayer = cameraService.getPreviewLayer()
        previewLayer.frame = view.bounds
        view.layer.addSublayer(previewLayer)

        context.coordinator.previewLayer = previewLayer

        return view
    }

    func updateUIView(_ uiView: UIView, context: Context) {
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
