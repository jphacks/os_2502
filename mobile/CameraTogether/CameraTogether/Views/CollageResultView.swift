import SwiftUI
import UIKit

struct CollageResultView: View {
    @Bindable var viewModel: CollageGroupViewModel
    let template: CollageTemplate
    let capturedImage: UIImage

    @State private var collageImage: UIImage?
    @State private var isGenerating = true
    @State private var errorMessage: String?
    @State private var showingSaveSuccess = false
    @Environment(\.colorScheme) private var colorScheme

    private var appColors: AppColors {
        AppColors(colorScheme: colorScheme)
    }

    var body: some View {
        ZStack {
            // 背景グラデーション
            LinearGradient(
                gradient: Gradient(colors: [
                    appColors.backgroundTop,
                    appColors.backgroundMiddle,
                    appColors.backgroundBottom
                ]),
                startPoint: .top,
                endPoint: .bottom
            )
            .ignoresSafeArea()

            VStack(spacing: 20) {
                // ヘッダー
                Text("コラージュ完成")
                    .font(.title)
                    .fontWeight(.bold)
                    .foregroundColor(appColors.textPrimary)
                    .padding(.top, 20)

                if isGenerating {
                    Spacer()
                    VStack(spacing: 16) {
                        ProgressView()
                            .scaleEffect(1.5)
                            .tint(appColors.textPrimary)
                        Text("コラージュを生成中...")
                            .foregroundColor(appColors.textSecondary)
                    }
                    Spacer()
                } else if let error = errorMessage {
                    Spacer()
                    VStack(spacing: 16) {
                        Image(systemName: "exclamationmark.triangle")
                            .font(.system(size: 50))
                            .foregroundColor(.red)
                        Text(error)
                            .foregroundColor(appColors.textPrimary)
                            .multilineTextAlignment(.center)
                            .padding()
                    }
                    Spacer()
                } else if let image = collageImage {
                    // コラージュ画像表示
                    ScrollView {
                        VStack(spacing: 20) {
                            Image(uiImage: image)
                                .resizable()
                                .aspectRatio(contentMode: .fit)
                                .cornerRadius(12)
                                .shadow(radius: 10)
                                .padding()

                            // アクションボタン
                            HStack(spacing: 16) {
                                Button {
                                    saveToPhotoLibrary(image: image)
                                } label: {
                                    VStack(spacing: 8) {
                                        Image(systemName: "square.and.arrow.down")
                                            .font(.title2)
                                        Text("保存")
                                            .font(.caption)
                                    }
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .padding()
                                    .background(Color.blue)
                                    .cornerRadius(12)
                                }

                                Button {
                                    shareImage(image: image)
                                } label: {
                                    VStack(spacing: 8) {
                                        Image(systemName: "square.and.arrow.up")
                                            .font(.title2)
                                        Text("共有")
                                            .font(.caption)
                                    }
                                    .foregroundColor(.white)
                                    .frame(maxWidth: .infinity)
                                    .padding()
                                    .background(Color.green)
                                    .cornerRadius(12)
                                }
                            }
                            .padding(.horizontal)
                        }
                    }
                }
            }
        }
        .alert("保存完了", isPresented: $showingSaveSuccess) {
            Button("OK", role: .cancel) {}
        } message: {
            Text("写真ライブラリに保存しました")
        }
        .task {
            await generateCollage()
        }
    }

    private func generateCollage() async {
        isGenerating = true
        errorMessage = nil

        do {
            // コラージュ生成
            let generator = CollageGenerator()
            let images = collectMemberImages()
            let image = try await generator.generateCollage(
                template: template,
                images: images
            )

            await MainActor.run {
                collageImage = image
                isGenerating = false
            }
        } catch {
            await MainActor.run {
                errorMessage = "コラージュの生成に失敗しました: \(error.localizedDescription)"
                isGenerating = false
            }
        }
    }

    private func collectMemberImages() -> [UIImage] {
        // 実際の実装では、各メンバーの撮影画像を取得
        // 現在は撮影した画像を複製して使用
        let memberCount = viewModel.currentGroup?.members.count ?? 1
        return Array(repeating: capturedImage, count: memberCount)
    }

    private func saveToPhotoLibrary(image: UIImage) {
        UIImageWriteToSavedPhotosAlbum(image, nil, nil, nil)
        showingSaveSuccess = true
    }

    private func shareImage(image: UIImage) {
        let activityVC = UIActivityViewController(
            activityItems: [image],
            applicationActivities: nil
        )

        if let windowScene = UIApplication.shared.connectedScenes.first as? UIWindowScene,
           let window = windowScene.windows.first,
           let rootVC = window.rootViewController {
            rootVC.present(activityVC, animated: true)
        }
    }
}

/// コラージュ生成器
class CollageGenerator {
    enum GeneratorError: Error, LocalizedError {
        case invalidTemplate
        case insufficientImages
        case renderingFailed

        var errorDescription: String? {
            switch self {
            case .invalidTemplate: return "テンプレートが無効です"
            case .insufficientImages: return "画像が不足しています"
            case .renderingFailed: return "レンダリングに失敗しました"
            }
        }
    }

    func generateCollage(template: CollageTemplate, images: [UIImage]) async throws -> UIImage {
        guard images.count >= template.photoCount else {
            throw GeneratorError.insufficientImages
        }

        // キャンバスサイズ（正方形）
        let canvasSize = CGSize(width: 1080, height: 1080)

        // viewBoxをパース
        let viewBoxComponents = template.viewBox.split(separator: " ").compactMap { Double($0) }
        guard viewBoxComponents.count == 4 else {
            throw GeneratorError.invalidTemplate
        }

        let viewBoxWidth = viewBoxComponents[2]
        let viewBoxHeight = viewBoxComponents[3]
        let scaleX = canvasSize.width / viewBoxWidth
        let scaleY = canvasSize.height / viewBoxHeight

        // 画像レンダリング
        let renderer = UIGraphicsImageRenderer(size: canvasSize)
        let collageImage = renderer.image { context in
            // 背景を白に
            UIColor.white.setFill()
            context.fill(CGRect(origin: .zero, size: canvasSize))

            // 各フレームに画像を配置
            for (index, frame) in template.frames.enumerated() {
                guard index < images.count else { break }

                let image = images[index]

                // SVGパスからCGPathを作成
                if let path = createPath(from: frame.path, scaleX: scaleX, scaleY: scaleY) {
                    // パスをクリッピングマスクとして使用
                    context.cgContext.saveGState()
                    context.cgContext.addPath(path)
                    context.cgContext.clip()

                    // パスの境界を取得
                    let bounds = path.boundingBox

                    // 画像を境界に合わせて描画
                    image.draw(in: bounds, blendMode: .normal, alpha: 1.0)

                    context.cgContext.restoreGState()

                    // フレームの境界線を描画
                    context.cgContext.addPath(path)
                    context.cgContext.setStrokeColor(UIColor.white.cgColor)
                    context.cgContext.setLineWidth(2)
                    context.cgContext.strokePath()
                }
            }
        }

        return collageImage
    }

    private func createPath(from pathString: String, scaleX: CGFloat, scaleY: CGFloat) -> CGPath? {
        let path = CGMutablePath()
        var currentPoint = CGPoint.zero

        let commands = pathString.uppercased()
        var index = commands.startIndex

        while index < commands.endIndex {
            let command = commands[index]
            index = commands.index(after: index)

            switch command {
            case "M":  // MoveTo
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 2) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: coords[1] * scaleY)
                    path.move(to: currentPoint)
                }
            case "L":  // LineTo
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 2) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: coords[1] * scaleY)
                    path.addLine(to: currentPoint)
                }
            case "H":  // Horizontal Line
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 1) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: currentPoint.y)
                    path.addLine(to: currentPoint)
                }
            case "V":  // Vertical Line
                if let coords = parseCoordinates(from: commands, startingAt: &index, count: 1) {
                    currentPoint = CGPoint(x: currentPoint.x, y: coords[0] * scaleY)
                    path.addLine(to: currentPoint)
                }
            case "Z":  // ClosePath
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

#Preview {
    let authManager = AuthenticationManager()
    let viewModel = CollageGroupViewModel(authManager: authManager)
    let _ = viewModel.createGroupLocal(type: .temporaryLocal, maxMembers: 2)

    let template = CollageTemplate(
        name: "2人用_縦分割",
        photoCount: 2,
        viewBox: "0 0 1 1",
        frames: [
            CollageTemplateFrame(id: 1, path: "M0.02 0.02H0.49V0.98H0.02V0.02Z"),
            CollageTemplateFrame(id: 2, path: "M0.51 0.02H0.98V0.98H0.51V0.02Z")
        ]
    )

    return NavigationStack {
        CollageResultView(
            viewModel: viewModel,
            template: template,
            capturedImage: UIImage(systemName: "photo")!
        )
    }
}
