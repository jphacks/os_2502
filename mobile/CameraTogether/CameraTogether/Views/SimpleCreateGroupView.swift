import SwiftUI

struct SimpleCreateGroupView: View {
    @State private var viewModel = CollageGroupViewModel()
    @State private var selectedGroupType: GroupType = .temporaryLocal
    @State private var showingWaitingRoom = false
    @Environment(\.appColors) var appColors

    var body: some View {
        ZStack {
            appColors.backgroundGradient
                .ignoresSafeArea()

            VStack(spacing: 30) {
                Text("グループ作成")
                    .font(.largeTitle)
                    .fontWeight(.bold)

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

                Button {
                    createGroup()
                } label: {
                    Text("グループを作成")
                        .font(.title2)
                        .fontWeight(.semibold)
                        .foregroundColor(.white)
                        .frame(maxWidth: .infinity)
                        .padding()
                        .background(Color.blue)
                        .cornerRadius(12)
                }
            }
            .padding()
        }
        .navigationDestination(isPresented: $showingWaitingRoom) {
            SimpleWaitingRoomView(viewModel: viewModel)
        }
    }

    private func createGroup() {
        viewModel.createGroup(type: selectedGroupType, maxMembers: 10)
        // モックデータとして他のメンバーを追加
        _ = viewModel.addMember(name: "太郎")
        _ = viewModel.addMember(name: "花子")
        showingWaitingRoom = true
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
        SimpleCreateGroupView()
    }
}
