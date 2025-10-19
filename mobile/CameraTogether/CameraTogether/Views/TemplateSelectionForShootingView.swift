import SwiftUI

/// 撮影用のテンプレート選択画面（オーナーのみ）- ランダムに自動選択
struct TemplateSelectionForShootingView: View {
    @Bindable var viewModel: CollageGroupViewModel
    let onTemplateSelected: (CollageTemplate, Int) -> Void

    @State private var templates: [CollageTemplate] = []
    @State private var isLoading = true
    @State private var errorMessage: String?
    @State private var selectedTemplate: CollageTemplate?
    @Environment(\.colorScheme) private var colorScheme
    @Environment(\.dismiss) private var dismiss

    private let templateService = TemplateAPIService.shared

    private var appColors: AppColors {
        AppColors(colorScheme: colorScheme)
    }

    var body: some View {
        NavigationStack {
            ZStack {
                LinearGradient(
                    gradient: Gradient(colors: [
                        appColors.backgroundTop,
                        appColors.backgroundMiddle,
                        appColors.backgroundBottom,
                    ]),
                    startPoint: .top,
                    endPoint: .bottom
                )
                .ignoresSafeArea()

                VStack(spacing: 20) {
                    if isLoading {
                        Spacer()
                        VStack(spacing: 16) {
                            ProgressView()
                                .scaleEffect(1.5)
                                .tint(appColors.textPrimary)
                            Text("テンプレートを読み込んでいます...")
                                .font(.caption)
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
                                .padding(.horizontal)

                            Button {
                                Task {
                                    await loadTemplates()
                                }
                            } label: {
                                HStack(spacing: 8) {
                                    Image(systemName: "arrow.clockwise")
                                    Text("再試行")
                                }
                                .foregroundColor(.white)
                                .padding(.horizontal, 24)
                                .padding(.vertical, 12)
                                .background(Color.blue)
                                .cornerRadius(12)
                            }
                        }
                        .padding()
                        Spacer()
                    } else if templates.isEmpty {
                        Spacer()
                        VStack(spacing: 16) {
                            Image(systemName: "photo.on.rectangle.angled")
                                .font(.system(size: 50))
                                .foregroundColor(appColors.textSecondary)
                            Text("利用可能なテンプレートがありません")
                                .foregroundColor(appColors.textSecondary)
                            Text(
                                "\(viewModel.currentGroup?.members.count ?? 0)人用のテンプレートが見つかりませんでした"
                            )
                            .font(.caption)
                            .foregroundColor(appColors.textSecondary)
                        }
                        Spacer()
                    } else {
                        Spacer()
                        VStack(spacing: 24) {
                            Image(systemName: "sparkles")
                                .font(.system(size: 60))
                                .foregroundColor(.yellow)

                            Text("テンプレートを選択中...")
                                .font(.title2)
                                .fontWeight(.bold)
                                .foregroundColor(appColors.textPrimary)

                            Text("ランダムにテンプレートを選択しています")
                                .font(.caption)
                                .foregroundColor(appColors.textSecondary)

                            ProgressView()
                                .scaleEffect(1.2)
                                .tint(appColors.textPrimary)
                        }
                        Spacer()
                    }
                }
            }
            .navigationBarTitleDisplayMode(.inline)
            .toolbar {
                ToolbarItem(placement: .navigationBarLeading) {
                    Button("キャンセル") {
                        dismiss()
                    }
                }
            }
        }
        .task {
            await loadTemplates()
        }
    }

    private func loadTemplates() async {
        print("loadTemplates: Starting...")
        isLoading = true
        errorMessage = nil

        do {
            let photoCount = viewModel.currentGroup?.members.count ?? 1
            print("loadTemplates: Fetching templates for \(photoCount) photos")

            templates = try await templateService.getTemplates(photoCount: photoCount)
            print("loadTemplates: Got \(templates.count) templates")

            isLoading = false

            // テンプレートが見つかった場合、ランダムに選択
            if !templates.isEmpty {
                // 少し待ってからランダム選択（UI効果のため）
                try? await Task.sleep(nanoseconds: 500_000_000)  // 0.5秒

                let randomTemplate = templates.randomElement()!
                print("Selected random template: \(randomTemplate.name)")

                // 現在のユーザーのインデックスを取得
                let members = viewModel.currentGroup?.members ?? []
                let myIndex = members.firstIndex(where: { $0.id == viewModel.currentUserId }) ?? 0

                print("My member index: \(myIndex)")

                await MainActor.run {
                    onTemplateSelected(randomTemplate, myIndex)
                }
            }
        } catch let error as APIError {
            print("loadTemplates API error: \(error)")
            errorMessage = "テンプレートの読み込みに失敗しました\n\(error.localizedDescription)"
            isLoading = false
        } catch {
            print("loadTemplates unknown error: \(error)")
            errorMessage = "テンプレートの読み込みに失敗しました\n\(error.localizedDescription)"
            isLoading = false
        }
    }
}
