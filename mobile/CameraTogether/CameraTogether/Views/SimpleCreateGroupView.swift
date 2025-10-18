import SwiftUI

struct SimpleCreateGroupView: View {
    let authManager: AuthenticationManager
    @State private var viewModel: CollageGroupViewModel?
    @State private var selectedGroupType: GroupType = .temporaryLocal
    @State private var showingWaitingRoom = false
    @State private var groupName: String = ""
    @Environment(\.appColors) var appColors
    @Environment(\.dismiss) var dismiss

    var body: some View {
        ZStack {
            appColors.backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 30) {
                Text("グループ作成")
                    .font(.largeTitle)
                    .fontWeight(.bold)

                // グループ名入力
                VStack(alignment: .leading, spacing: 8) {
                    Text("グループ名")
                        .font(.headline)
                        .foregroundColor(.primary)
                    TextField("例: 友達グループ", text: $groupName)
                        .textFieldStyle(RoundedBorderTextFieldStyle())
                        .padding(.horizontal, 4)
                }
                .padding(.horizontal, 16)

                VStack(spacing: 20) {
                    GroupTypeButton(
                        title: "ローカル",
                        description: "近くにいる友達と簡単にグループを作成",
                        icon: "location.fill",
                        isSelected: selectedGroupType == .temporaryLocal
                    ) {
                        selectedGroupType = .temporaryLocal
                    }

                    GroupTypeButton(
                        title: "グローバル",
                        description: "インターネットを通じて友達とグループを作成",
                        icon: "globe",
                        isSelected: selectedGroupType == .temporaryGlobal
                    ) {
                        selectedGroupType = .temporaryGlobal
                    }

                    GroupTypeButton(
                        title: "固定",
                        description: "いつでも参加できる固定グループを作成",
                        icon: "lock.fill",
                        isSelected: selectedGroupType == .fixed
                    ) {
                        selectedGroupType = .fixed
                    }
                }
                .padding(.horizontal, 16)

                Spacer()

                // エラーメッセージ
                if let errorMessage = viewModel?.errorMessage {
                    Text(errorMessage)
                        .font(.caption)
                        .foregroundColor(.red)
                        .multilineTextAlignment(.center)
                        .padding(.horizontal)
                }

                Button {
                    Task {
                        await createGroup()
                    }
                } label: {
                    if viewModel?.isLoading == true {
                        ProgressView()
                            .progressViewStyle(CircularProgressViewStyle(tint: .white))
                            .frame(maxWidth: .infinity)
                            .frame(height: 56)
                    } else {
                        Text("グループを作成")
                            .font(.title2)
                            .fontWeight(.semibold)
                            .foregroundColor(.white)
                            .frame(maxWidth: .infinity)
                            .frame(height: 56)
                    }
                }
                .background(groupName.isEmpty ? Color.gray : Color.blue)
                .cornerRadius(12)
                .disabled(groupName.isEmpty || viewModel?.isLoading == true)
                .padding(.horizontal)
            }
            .padding()
        }
        .navigationDestination(isPresented: $showingWaitingRoom) {
            if let viewModel = viewModel {
                SimpleWaitingRoomView(viewModel: viewModel)
            }
        }
        .onAppear {
            if viewModel == nil {
                viewModel = CollageGroupViewModel(authManager: authManager)
            }
        }
    }

    private func createGroup() async {
        let name = groupName.trimmingCharacters(in: .whitespacesAndNewlines)
        guard !name.isEmpty else { return }
        guard let vm = viewModel else { return }

        // API経由でグループ作成
        await vm.createGroup(
            type: selectedGroupType,
            name: name,
            maxMembers: 10
        )

        // 作成成功したら待機室へ遷移
        if vm.currentGroup != nil {
            showingWaitingRoom = true
        }
    }
}

struct GroupTypeButton: View {
    let title: String
    let description: String
    let icon: String
    let isSelected: Bool
    let action: () -> Void

    var body: some View {
        Button(action: action) {
            HStack(spacing: 16) {
                Image(systemName: icon)
                    .font(.title)
                    .foregroundColor(isSelected ? .blue : .gray)
                    .frame(width: 50)

                VStack(alignment: .leading, spacing: 4) {
                    Text(title)
                        .font(.headline)
                        .foregroundColor(.primary)
                    Text(description)
                        .font(.caption)
                        .foregroundColor(.gray)
                }

                Spacer()

                if isSelected {
                    Image(systemName: "checkmark.circle.fill")
                        .foregroundColor(.blue)
                        .font(.title2)
                }
            }
            .padding()
            .background(isSelected ? Color.blue.opacity(0.1) : Color.gray.opacity(0.1))
            .cornerRadius(12)
            .overlay(
                RoundedRectangle(cornerRadius: 12)
                    .stroke(isSelected ? Color.blue : Color.clear, lineWidth: 2)
            )
        }
    }
}

#Preview {
    NavigationStack {
        SimpleCreateGroupView(authManager: AuthenticationManager())
    }
}
