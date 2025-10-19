import SwiftUI

struct TemplateSelectionView: View {
    @Bindable var viewModel: CollageGroupViewModel
    let capturedImage: UIImage

    @State private var templates: [CollageTemplate] = []
    @State private var isLoading = true
    @State private var errorMessage: String?
    @State private var selectedTemplate: CollageTemplate?
    @State private var showingCollageResult = false
    @Environment(\.colorScheme) private var colorScheme

    private var appColors: AppColors {
        AppColors(colorScheme: colorScheme)
    }

    private let templateService = TemplateAPIService.shared

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
                Text("コラージュテンプレート選択")
                    .font(.title)
                    .fontWeight(.bold)
                    .foregroundColor(appColors.textPrimary)
                    .padding(.top, 20)

                if isLoading {
                    Spacer()
                    ProgressView()
                        .scaleEffect(1.5)
                        .tint(appColors.textPrimary)
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
                        Button("再読み込み") {
                            Task {
                                await loadTemplates()
                            }
                        }
                        .buttonStyle(.bordered)
                    }
                    .padding()
                    Spacer()
                } else if templates.isEmpty {
                    Spacer()
                    Text("利用可能なテンプレートがありません")
                        .foregroundColor(appColors.textSecondary)
                    Spacer()
                } else {
                    ScrollView {
                        LazyVGrid(columns: [
                            GridItem(.flexible()),
                            GridItem(.flexible())
                        ], spacing: 16) {
                            ForEach(templates) { template in
                                TemplateCard(
                                    template: template,
                                    isSelected: selectedTemplate?.id == template.id
                                ) {
                                    selectedTemplate = template
                                }
                            }
                        }
                        .padding(.horizontal)
                    }

                    if selectedTemplate != nil {
                        Button {
                            showingCollageResult = true
                        } label: {
                            HStack(spacing: 12) {
                                Image(systemName: "wand.and.stars")
                                    .font(.title3)
                                Text("コラージュ作成")
                                    .font(.headline)
                            }
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .padding(.vertical, 16)
                            .background(
                                LinearGradient(
                                    colors: [.blue, .blue.opacity(0.8)],
                                    startPoint: .leading,
                                    endPoint: .trailing
                                )
                            )
                            .cornerRadius(16)
                        }
                        .padding(.horizontal, 24)
                        .padding(.bottom, 16)
                    }
                }
            }
        }
        .task {
            await loadTemplates()
        }
        .navigationDestination(isPresented: $showingCollageResult) {
            if let template = selectedTemplate {
                CollageResultView(
                    viewModel: viewModel,
                    template: template,
                    capturedImage: capturedImage
                )
            }
        }
    }

    private func loadTemplates() async {
        isLoading = true
        errorMessage = nil

        do {
            // グループのメンバー数に応じてテンプレートを取得
            let photoCount = viewModel.currentGroup?.members.count ?? 1
            templates = try await templateService.getTemplates(photoCount: photoCount)
            isLoading = false
        } catch {
            errorMessage = "テンプレートの読み込みに失敗しました: \(error.localizedDescription)"
            isLoading = false
        }
    }
}

struct TemplateCard: View {
    let template: CollageTemplate
    let isSelected: Bool
    let onTap: () -> Void

    var body: some View {
        Button(action: onTap) {
            VStack(spacing: 8) {
                // テンプレートプレビュー
                ZStack {
                    RoundedRectangle(cornerRadius: 12)
                        .fill(Color.white.opacity(0.1))
                        .aspectRatio(1, contentMode: .fit)

                    // SVGパスのプレビュー（簡易表示）
                    TemplatePreview(template: template)
                        .padding(8)
                }
                .overlay(
                    RoundedRectangle(cornerRadius: 12)
                        .stroke(isSelected ? Color.blue : Color.white.opacity(0.3), lineWidth: isSelected ? 3 : 1)
                )

                Text(template.name)
                    .font(.caption)
                    .foregroundColor(.white)
                    .lineLimit(2)
                    .multilineTextAlignment(.center)
                    .fixedSize(horizontal: false, vertical: true)
            }
        }
        .buttonStyle(.plain)
    }
}

struct TemplatePreview: View {
    let template: CollageTemplate

    var body: some View {
        GeometryReader { geometry in
            ZStack {
                ForEach(template.frames) { frame in
                    FrameShape(pathString: frame.path, viewBox: template.viewBox)
                        .stroke(Color.white, lineWidth: 2)
                }
            }
        }
    }
}

struct FrameShape: Shape {
    let pathString: String
    let viewBox: String

    func path(in rect: CGRect) -> Path {
        // viewBoxをパース
        let components = viewBox.split(separator: " ").compactMap { Double($0) }
        guard components.count == 4 else {
            return Path()
        }

        let viewBoxWidth = components[2]
        let viewBoxHeight = components[3]

        // スケール計算
        let scaleX = rect.width / viewBoxWidth
        let scaleY = rect.height / viewBoxHeight

        var path = Path()

        // SVGパスコマンドを簡易パース
        let commands = pathString.uppercased()
        var currentPoint = CGPoint.zero
        var commandIndex = commands.startIndex

        while commandIndex < commands.endIndex {
            let command = commands[commandIndex]
            commandIndex = commands.index(after: commandIndex)

            switch command {
            case "M":  // MoveTo
                if let coords = parseCoordinates(from: commands, startingAt: &commandIndex, count: 2) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: coords[1] * scaleY)
                    path.move(to: currentPoint)
                }
            case "L":  // LineTo
                if let coords = parseCoordinates(from: commands, startingAt: &commandIndex, count: 2) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: coords[1] * scaleY)
                    path.addLine(to: currentPoint)
                }
            case "H":  // Horizontal Line
                if let coords = parseCoordinates(from: commands, startingAt: &commandIndex, count: 1) {
                    currentPoint = CGPoint(x: coords[0] * scaleX, y: currentPoint.y)
                    path.addLine(to: currentPoint)
                }
            case "V":  // Vertical Line
                if let coords = parseCoordinates(from: commands, startingAt: &commandIndex, count: 1) {
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

            if char.isNumber || char == "." {
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
    let _ = viewModel.createGroupLocal(type: .temporaryLocal, maxMembers: 3)

    return NavigationStack {
        TemplateSelectionView(
            viewModel: viewModel,
            capturedImage: UIImage(systemName: "photo")!
        )
    }
}
