import CoreImage.CIFilterBuiltins
import SwiftUI

struct SimpleWaitingRoomView: View {
    @Bindable var viewModel: CollageGroupViewModel
    @State private var showingCountdown = false
    @State private var showAddMemberSheet = false
    @State private var showQRCodeSheet = false
    @State private var showFriendList = false
    @State private var selectedTemplate: CollageTemplate?
    @State private var myFrameIndex: Int = 0
    @State private var isLoadingTemplate = false
    @State private var templateError: String?
    @Environment(\.dismiss) private var dismiss
    @Environment(\.colorScheme) var colorScheme
    @Environment(\.appColors) var appColors

    var body: some View {
        ZStack {
            appColors.backgroundGradient
                .ignoresSafeArea()

            if let group = viewModel.currentGroup {
                VStack(spacing: 0) {
                    ScrollView {
                        VStack(spacing: 24) {
                            Spacer()
                                .frame(height: 8)

                            memberListSection(group: group)

                            Spacer()
                                .frame(height: 40)
                        }
                        .padding(.horizontal, 24)
                    }

                    if let currentMember = group.members.first(where: {
                        $0.id == viewModel.currentUserId
                    }) {
                        readyButtonSection(currentMember: currentMember)
                    }
                }
            }
        }
        .navigationTitle(groupTypeText(type: viewModel.currentGroup?.type ?? .temporaryLocal))
        .navigationBarTitleDisplayMode(.inline)
        .toolbar {
            ToolbarItem(placement: .navigationBarLeading) {
                Button {
                    dismiss()
                } label: {
                    HStack(spacing: 4) {
                        Image(systemName: "xmark")
                        Text("閉じる")
                    }
                    .foregroundColor(appColors.textPrimary)
                }
            }

            ToolbarItem(placement: .navigationBarTrailing) {
                if viewModel.isOwner {
                    Button {
                        showAddMemberSheet = true
                    } label: {
                        Image(systemName: "plus.circle.fill")
                            .font(.title3)
                            .foregroundColor(appColors.textPrimary)
                    }
                }
            }
        }
        .navigationBarBackButtonHidden(true)
        .toolbarBackground(.visible, for: .navigationBar)
        .toolbarBackground(Color.clear, for: .navigationBar)
        .navigationDestination(isPresented: $showingCountdown) {
            if let template = selectedTemplate {
                CountdownWithGuideView(
                    viewModel: viewModel,
                    template: template,
                    myFrameIndex: myFrameIndex
                )
            } else {
                let _ = print("Navigation: selectedTemplate is nil!")
                Text("テンプレートが選択されていません")
                    .foregroundColor(.red)
            }
        }
        .onChange(of: showingCountdown) {
            if !showingCountdown {
                print("showingCountdown was set to false! Navigation dismissed.")
            }
        }
        .sheet(isPresented: $showAddMemberSheet) {
            AddMemberSheetView(
                groupType: viewModel.currentGroup?.type ?? .temporaryLocal,
                onShowQR: {
                    showAddMemberSheet = false
                    showQRCodeSheet = true
                },
                onFriendSelect: {
                    showAddMemberSheet = false
                    showFriendList = true
                }
            )
            .presentationDetents([.height(260)])
        }
        .sheet(isPresented: $showQRCodeSheet) {
            if let group = viewModel.currentGroup {
                GroupQRCodeView(group: group)
            }
        }
        .sheet(isPresented: $showFriendList) {
            FriendSelectView { friendName in
                _ = viewModel.addMember(name: friendName)
                showFriendList = false
            }
        }
        .onAppear {
            viewModel.startMemberPolling()
        }
        .onDisappear {
            viewModel.stopMemberPolling()
        }
        .onChange(of: viewModel.currentGroup?.status) {
            let newStatus = viewModel.currentGroup?.status

            // ステータスがcountdownに変わったら、参加者も自動的に撮影画面に遷移
            if newStatus == .countdown && !viewModel.isOwner && !showingCountdown {
                Task {
                    await loadRandomTemplateForParticipant()
                }
            }
        }
    }

    @ViewBuilder
    private func memberListSection(group: CollageGroup) -> some View {
        VStack(spacing: 16) {
            HStack {
                HStack(spacing: 8) {
                    ZStack {
                        Circle()
                            .fill(
                                LinearGradient(
                                    colors: [.blue.opacity(0.6), .cyan.opacity(0.4)],
                                    startPoint: .topLeading,
                                    endPoint: .bottomTrailing
                                )
                            )
                            .frame(width: 28, height: 28)
                        Image(systemName: "person.2.fill")
                            .font(.system(size: 14))
                            .foregroundColor(.white)
                    }
                    Text("参加メンバー")
                        .font(.title3)
                        .fontWeight(.bold)
                        .foregroundColor(appColors.textPrimary)
                }
                Spacer()
                Text("\(group.members.count)人")
                    .font(.subheadline)
                    .foregroundColor(appColors.textSecondary)
            }

            VStack(spacing: 12) {
                ForEach(group.members) { member in
                    MemberCardView(member: member)
                }
            }
        }
    }

    @ViewBuilder
    private func readyButtonSection(currentMember: CollageGroupMember) -> some View {
        VStack(spacing: 0) {
            Divider()
                .background(Color.white.opacity(0.1))

            if let group = viewModel.currentGroup {
                let _ = print(
                    "Button section - finalized: \(group.isFinalized), allReady: \(group.allMembersReady), isOwner: \(viewModel.isOwner), memberCount: \(group.members.count)"
                )

                if !group.isFinalized && viewModel.isOwner && group.members.count == 1 {
                    let _ = print("Showing: メンバー招待メッセージ")
                    // グループ未確定 & オーナー & 1人のみ: 招待メッセージ
                    VStack(spacing: 12) {
                        Image(systemName: "person.badge.plus")
                            .font(.largeTitle)
                            .foregroundColor(appColors.textSecondary)
                        Text("メンバーを招待してください")
                            .font(.headline)
                            .foregroundColor(appColors.textPrimary)
                        Text("右上の + ボタンからメンバーを追加できます")
                            .font(.caption)
                            .foregroundColor(appColors.textSecondary)
                            .multilineTextAlignment(.center)
                    }
                    .padding(.vertical, 24)
                } else if !group.isFinalized && viewModel.isOwner && group.members.count > 1 {
                    let _ = print("Showing: グループ確定ボタン")
                    // グループ未確定 & オーナー & 2人以上: グループ確定ボタン
                    Button {
                        print("グループ確定ボタンがタップされました")
                        Task {
                            await viewModel.finalizeGroup()
                        }
                    } label: {
                        HStack(spacing: 12) {
                            Image(systemName: "checkmark.circle.fill")
                                .font(.title3)
                            Text("このメンバーでコラージュを作る (\(group.members.count)人)")
                                .font(.headline)
                        }
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 16)
                        .background(
                            LinearGradient(
                                colors: [.green, .green.opacity(0.8)],
                                startPoint: .leading,
                                endPoint: .trailing
                            )
                        )
                        .cornerRadius(16)
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 16)
                } else if group.isFinalized && group.allMembersReady && viewModel.isOwner {
                    let _ = print("Showing: 撮影開始ボタン (isLoadingTemplate: \(isLoadingTemplate))")
                    // グループ確定済 & 全員準備完了 & オーナー: 撮影開始ボタン
                    VStack(spacing: 12) {
                        Button {
                            print("撮影開始ボタンがタップされました")
                            Task {
                                await loadRandomTemplateAndStart()
                            }
                        } label: {
                            HStack(spacing: 12) {
                                if isLoadingTemplate {
                                    ProgressView()
                                        .progressViewStyle(CircularProgressViewStyle(tint: .white))
                                } else {
                                    Image(systemName: "camera.fill")
                                        .font(.title3)
                                }
                                Text(isLoadingTemplate ? "準備中..." : "撮影開始")
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
                        .disabled(isLoadingTemplate)

                        // エラー表示
                        if let error = templateError {
                            Text(error)
                                .font(.caption)
                                .foregroundColor(.red)
                                .multilineTextAlignment(.center)
                        }
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 16)
                } else if group.isFinalized && !currentMember.isReady {
                    let _ = print("Showing: 準備完了ボタン")
                    // グループ確定済 & 未準備: 準備完了ボタン
                    Button {
                        Task {
                            await viewModel.markReadyWithAPI()
                        }
                    } label: {
                        HStack(spacing: 12) {
                            Image(systemName: "checkmark.circle.fill")
                                .font(.title3)
                            Text("準備完了")
                                .font(.headline)
                        }
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .padding(.vertical, 16)
                        .background(
                            LinearGradient(
                                colors: [.orange, .orange.opacity(0.8)],
                                startPoint: .leading,
                                endPoint: .trailing
                            )
                        )
                        .cornerRadius(16)
                    }
                    .padding(.horizontal, 24)
                    .padding(.vertical, 16)
                } else if !group.isFinalized && !viewModel.isOwner {
                    let _ = print("Showing: 待機メッセージ（オーナーの確定待ち）")
                    // グループ未確定 & 非オーナー: 待機メッセージ
                    VStack(spacing: 8) {
                        Image(systemName: "clock")
                            .font(.title2)
                            .foregroundColor(appColors.textSecondary)
                        Text("オーナーがメンバーを確定するまで待機中...")
                            .font(.caption)
                            .foregroundColor(appColors.textSecondary)
                            .multilineTextAlignment(.center)
                    }
                    .padding(.vertical, 16)
                } else if group.isFinalized && currentMember.isReady {
                    let _ = print("Showing: 準備完了済み表示")
                    // 準備完了済み: 待機中
                    VStack(spacing: 8) {
                        HStack(spacing: 12) {
                            Image(systemName: "checkmark.circle.fill")
                                .font(.title3)
                                .foregroundColor(.green)
                            Text("準備完了")
                                .font(.headline)
                                .foregroundColor(.green)
                        }

                        if let group = viewModel.currentGroup, group.allMembersReady {
                            Text("撮影開始を待っています")
                                .font(.caption)
                                .foregroundColor(appColors.textSecondary)
                        } else {
                            Text("他のメンバーの準備完了を待っています")
                                .font(.caption)
                                .foregroundColor(appColors.textSecondary)
                        }
                    }
                    .padding(.vertical, 16)
                } else {
                    let _ = print("No button condition matched!")
                    Text("状態エラー")
                        .foregroundColor(.red)
                        .padding()
                }
            } else {
                let _ = print("currentGroup is nil")
            }
        }
        .background(
            appColors.backgroundGradient
                .opacity(0.95)
        )
    }

    private func groupTypeText(type: GroupType) -> String {
        switch type {
        case .temporaryLocal:
            return "ローカル"
        case .temporaryGlobal:
            return "グローバル"
        case .fixed:
            return "固定グループ"
        }
    }

    /// ランダムにテンプレートを選択して撮影開始
    private func loadRandomTemplateAndStart() async {
        await MainActor.run {
            isLoadingTemplate = true
            templateError = nil
        }

        guard let group = viewModel.currentGroup else {
            print("Group not found")
            await MainActor.run {
                isLoadingTemplate = false
                templateError = "グループ情報が見つかりません"
            }
            return
        }

        let photoCount = group.members.count
        print("Loading templates for \(photoCount) photos...")

        do {
            // テンプレートAPIから取得
            let templateService = TemplateAPIService.shared
            let templates = try await templateService.getTemplates(photoCount: photoCount)

            print("Got \(templates.count) templates")

            guard !templates.isEmpty else {
                print("No templates found for \(photoCount) photos")
                await MainActor.run {
                    isLoadingTemplate = false
                    templateError = "\(photoCount)人用のテンプレートが見つかりません"
                }
                return
            }

            // ランダムに1つ選択
            let randomTemplate = templates.randomElement()!
            print("Selected random template: \(randomTemplate.name), id: \(randomTemplate.id)")

            // 自分のフレームインデックスを取得
            let members = group.members
            let myIndex = members.firstIndex(where: { $0.id == viewModel.currentUserId }) ?? 0
            print("My frame index: \(myIndex)")

            // 状態を更新
            await MainActor.run {
                selectedTemplate = randomTemplate
                myFrameIndex = myIndex
            }

            // API呼び出し（テンプレートIDを渡す）
            let success = await viewModel.startCountdownWithAPI(templateId: randomTemplate.id)

            await MainActor.run {
                isLoadingTemplate = false
                if success {
                    if viewModel.currentGroup?.scheduledCaptureTime != nil {
                        showingCountdown = true
                    } else {
                        templateError = "撮影時刻の取得に失敗しました"
                    }
                } else {
                    templateError = "撮影開始に失敗しました"
                }
            }
        } catch {
            print("Failed to load templates: \(error.localizedDescription)")
            await MainActor.run {
                isLoadingTemplate = false
                templateError = "テンプレート読み込みエラー: \(error.localizedDescription)"
            }
        }
    }

    /// 参加者用: サーバーから取得したテンプレートIDを使ってテンプレートをロードし遷移
    private func loadRandomTemplateForParticipant() async {
        await MainActor.run {
            isLoadingTemplate = true
            templateError = nil
        }

        guard let group = viewModel.currentGroup else {
            print("Group not found")
            await MainActor.run {
                isLoadingTemplate = false
            }
            return
        }

        // サーバーから取得したテンプレートIDを確認
        guard let templateId = group.templateId else {
            print("[Participant] Template ID not found in group")
            await MainActor.run {
                isLoadingTemplate = false
                templateError = "テンプレートIDが見つかりません"
            }
            return
        }

        print("[Participant] Using template ID from server: \(templateId)")

        do {
            let templateService = TemplateAPIService.shared
            // 指定されたテンプレートIDでテンプレートを取得
            let template = try await templateService.getTemplate(id: templateId)

            print("[Participant] Loaded template: \(template.name)")

            let members = group.members
            let myIndex = members.firstIndex(where: { $0.id == viewModel.currentUserId }) ?? 0
            print("[Participant] My frame index: \(myIndex)")

            await MainActor.run {
                selectedTemplate = template
                myFrameIndex = myIndex
                isLoadingTemplate = false

                if viewModel.currentGroup?.scheduledCaptureTime != nil {
                    showingCountdown = true
                } else {
                    templateError = "撮影時刻の取得に失敗しました"
                }
            }
        } catch {
            print("[Participant] Failed to load template: \(error.localizedDescription)")
            await MainActor.run {
                isLoadingTemplate = false
                templateError = "テンプレート読み込みエラー: \(error.localizedDescription)"
            }
        }
    }

}

#Preview {
    let authManager = AuthenticationManager()
    let viewModel = CollageGroupViewModel(authManager: authManager)
    let _ = viewModel.createGroupLocal(type: .temporaryLocal, maxMembers: 5)

    return NavigationStack {
        SimpleWaitingRoomView(viewModel: viewModel)
    }
}
