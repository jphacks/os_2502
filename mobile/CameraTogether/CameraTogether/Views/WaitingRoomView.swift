import SwiftUI

struct WaitingRoomView: View {
    @Bindable var viewModel: CollageGroupViewModel
    @State private var showingCountdown = false

    var body: some View {
        VStack(spacing: 30) {
            if let group = viewModel.currentGroup {
                switch group.status {
                case .recruiting:
                    recruitingView(group: group)
                case .readyCheck:
                    readyCheckView(group: group)
                case .countdown, .completed:
                    EmptyView()
                }
            }
        }
        .padding()
        .navigationBarBackButtonHidden(true)
        .navigationDestination(isPresented: $showingCountdown) {
            CountdownView(viewModel: viewModel)
        }
        .onChange(of: viewModel.currentGroup?.status) { _, newStatus in
            if newStatus == .countdown {
                showingCountdown = true
            }
        }
    }

    @ViewBuilder
    private func recruitingView(group: CollageGroup) -> some View {
        VStack(spacing: 30) {
            Text("メンバー募集中")
                .font(.largeTitle)
                .fontWeight(.bold)

            VStack(spacing: 16) {
                HStack {
                    Text("招待コード")
                        .font(.headline)
                    Spacer()
                }

                HStack {
                    Text(group.inviteCode)
                        .font(.system(size: 36, weight: .bold, design: .monospaced))
                        .foregroundColor(.blue)

                    Button {
                        UIPasteboard.general.string = group.inviteCode
                    } label: {
                        Image(systemName: "doc.on.doc")
                            .font(.title2)
                    }
                }
                .padding()
                .background(Color.blue.opacity(0.1))
                .cornerRadius(12)

                Text("このコードを共有してメンバーを招待")
                    .font(.caption)
                    .foregroundColor(.gray)
            }
            .padding()
            .background(Color.gray.opacity(0.1))
            .cornerRadius(12)

            VStack(spacing: 12) {
                HStack {
                    Text("参加メンバー")
                        .font(.headline)
                    Spacer()
                    Text("\(group.members.count) / \(group.maxMembers)")
                        .foregroundColor(.gray)
                }

                ScrollView {
                    VStack(spacing: 8) {
                        ForEach(group.members) { member in
                            HStack {
                                Image(systemName: "person.circle.fill")
                                    .font(.title2)
                                    .foregroundColor(.blue)
                                Text(member.name)
                                    .font(.body)
                                Spacer()
                                if member.id == group.ownerId {
                                    Text("オーナー")
                                        .font(.caption)
                                        .padding(.horizontal, 8)
                                        .padding(.vertical, 4)
                                        .background(Color.orange.opacity(0.2))
                                        .cornerRadius(8)
                                }
                            }
                            .padding()
                            .background(Color.white)
                            .cornerRadius(8)
                        }
                    }
                }
                .frame(maxHeight: 300)
            }
            .padding()
            .background(Color.gray.opacity(0.1))
            .cornerRadius(12)

            Spacer()

            if viewModel.isOwner {
                Button {
                    viewModel.startReadyCheck()
                } label: {
                    Text("このメンバーでコラージュをする")
                        .font(.title2)
                        .fontWeight(.semibold)
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(group.members.count >= 2 ? Color.blue : Color.gray)
                        .cornerRadius(12)
                }
                .disabled(group.members.count < 2)
            }
        }
    }

    @ViewBuilder
    private func readyCheckView(group: CollageGroup) -> some View {
        VStack(spacing: 30) {
            Text("準備確認")
                .font(.largeTitle)
                .fontWeight(.bold)

            Text("全員が準備完了するとカウントダウンが始まります")
                .font(.body)
                .foregroundColor(.gray)
                .multilineTextAlignment(.center)

            VStack(spacing: 12) {
                HStack {
                    Text("メンバー準備状況")
                        .font(.headline)
                    Spacer()
                }

                ScrollView {
                    VStack(spacing: 8) {
                        ForEach(group.members) { member in
                            HStack {
                                Image(systemName: "person.circle.fill")
                                    .font(.title2)
                                    .foregroundColor(.blue)
                                Text(member.name)
                                    .font(.body)
                                Spacer()
                                if member.isReady {
                                    Image(systemName: "checkmark.circle.fill")
                                        .foregroundColor(.green)
                                        .font(.title2)
                                } else {
                                    Image(systemName: "hourglass")
                                        .foregroundColor(.orange)
                                        .font(.title2)
                                }
                            }
                            .padding()
                            .background(Color.white)
                            .cornerRadius(8)
                        }
                    }
                }
                .frame(maxHeight: 300)
            }
            .padding()
            .background(Color.gray.opacity(0.1))
            .cornerRadius(12)

            Spacer()

            if let currentMember = group.members.first(where: { $0.id == viewModel.currentUserId })
            {
                if !currentMember.isReady {
                    Button {
                        viewModel.markReady()
                    } label: {
                        Text("準備完了")
                            .font(.title2)
                            .fontWeight(.semibold)
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .padding()
                            .background(Color.green)
                            .cornerRadius(12)
                    }
                } else {
                    Text("準備完了しました")
                        .font(.title2)
                        .fontWeight(.semibold)
                        .foregroundColor(.green)
                        .padding()
                }
            }
        }
    }
}

#Preview {
    @Previewable @State var viewModel = CollageGroupViewModel()
    let _ = viewModel.createGroup(type: .temporaryLocal, maxMembers: 5)

    NavigationStack {
        WaitingRoomView(viewModel: viewModel)
    }
}
