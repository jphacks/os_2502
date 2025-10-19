import SwiftUI

struct CreateGroupView: View {
    let authManager: AuthenticationManager
    @State private var viewModel: CollageGroupViewModel?
    @State private var selectedGroupType: GroupType = .temporaryLocal
    @State private var maxMembers: Int = 10
    @State private var showingWaitingRoom = false

    var body: some View {
        VStack(spacing: 30) {
            Text("グループ作成")
                .font(.largeTitle)
                .fontWeight(.bold)

            VStack(spacing: 20) {
                HStack {
                    Text("グループタイプ")
                        .font(.headline)
                    Spacer()
                }

                Picker("グループタイプ", selection: $selectedGroupType) {
                    Text("ローカル").tag(GroupType.temporaryLocal)
                    Text("グローバル").tag(GroupType.temporaryGlobal)
                    Text("固定").tag(GroupType.fixed)
                }
                .pickerStyle(.segmented)

                HStack {
                    Text("最大メンバー数")
                        .font(.headline)
                    Spacer()
                    Text("\(maxMembers)人")
                        .foregroundColor(.gray)
                }

                Stepper("最大メンバー数", value: $maxMembers, in: 2...20)
                    .labelsHidden()
            }
            .padding()
            .background(Color.gray.opacity(0.1))
            .cornerRadius(12)

            Spacer()

            Button {
                Task {
                    if let vm = viewModel {
                        await vm.createGroup(
                            type: selectedGroupType,
                            name: "\(vm.currentUserName)のグループ",
                            maxMembers: maxMembers
                        )
                        await MainActor.run {
                            showingWaitingRoom = true
                        }
                    }
                }
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
        .navigationDestination(isPresented: $showingWaitingRoom) {
            if let viewModel = viewModel {
                WaitingRoomView(viewModel: viewModel)
            }
        }
        .onAppear {
            if viewModel == nil {
                viewModel = CollageGroupViewModel(authManager: authManager)
            }
        }
    }
}

#Preview {
    NavigationStack {
        CreateGroupView(authManager: AuthenticationManager())
    }
}
