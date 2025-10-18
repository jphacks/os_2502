import CoreImage.CIFilterBuiltins
import SwiftUI

struct SimpleWaitingRoomView: View {
    @Bindable var viewModel: CollageGroupViewModel
    @State private var showingCountdown = false
    @State private var showAddMemberSheet = false
    @State private var showQRCodeSheet = false
    @State private var showFriendList = false
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
            ToolbarItem(placement: .navigationBarTrailing) {
                Button {
                    showAddMemberSheet = true
                } label: {
                    Image(systemName: "plus.circle.fill")
                        .font(.title3)
                        .foregroundColor(appColors.textPrimary)
                }
            }
        }
        .toolbarBackground(.visible, for: .navigationBar)
        .toolbarBackground(Color.clear, for: .navigationBar)
        .navigationDestination(isPresented: $showingCountdown) {
            CountdownView(viewModel: viewModel)
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
            .presentationDetents([
                .height(viewModel.currentGroup?.type == .temporaryLocal ? 260 : 180)
            ])
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

            if !currentMember.isReady {
                Button {
                    viewModel.markReady()
                    checkAllReady()
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
                            colors: [.green, .green.opacity(0.8)],
                            startPoint: .leading,
                            endPoint: .trailing
                        )
                    )
                    .cornerRadius(16)
                }
                .padding(.horizontal, 24)
                .padding(.vertical, 16)
            } else {
                VStack(spacing: 8) {
                    HStack(spacing: 12) {
                        Image(systemName: "checkmark.circle.fill")
                            .font(.title3)
                            .foregroundColor(.green)
                        Text("準備完了")
                            .font(.headline)
                            .foregroundColor(.green)
                    }

                    Text("他のメンバーの準備完了を待っています")
                        .font(.caption)
                        .foregroundColor(appColors.textSecondary)
                }
                .padding(.vertical, 16)
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

    private func checkAllReady() {
        if let group = viewModel.currentGroup, group.allMembersReady {
            DispatchQueue.main.asyncAfter(deadline: .now() + 0.5) {
                showingCountdown = true
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
